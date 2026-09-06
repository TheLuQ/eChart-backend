package firestore

import (
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

const sheetStorageBaseURLEnv = "SHEET_STORAGE_BASE_URL"
const sheetStorageBucketEnv = "SHEET_STORAGE_BUCKET"

var groupKeySeparatorRegex = regexp.MustCompile(`[^a-z0-9]+`)

type Instrument struct {
	Name    string `firestore:"instrument_name_en" json:"instrument_name_en"`
	NamePol string `firestore:"instrument_name_pol" json:"instrument_name_pol"`
	Voice   string `firestore:"voice,omitempty" json:"voice,omitempty"`
	Key     string `firestore:"key,omitempty" json:"key,omitempty"`
}

type Sheet struct {
	Instrument
	Id       string `firestore:"object_name" json:"file_id"`
	Title    string `firestore:"-" json:"title"`
	Band     string `firestore:"-" json:"band,omitempty"`
	FileName string `firestore:"file_name,omitempty" json:"file_name,omitempty"`
	Url      string `firestore:"-" json:"url,omitempty"`
}

type SheetGroup struct {
	ID          string  `firestore:"-" json:"id"`
	Title       string  `firestore:"title,omitempty" json:"title"`
	Band        string  `firestore:"band,omitempty" json:"band,omitempty"`
	GroupKey    string  `firestore:"group_key,omitempty" json:"group_key,omitempty"`
	LastUpdated string  `firestore:"last_updated" json:"last_updated,omitempty"`
	ParentPath  string  `firestore:"parent_path,omitempty" json:"parent_path,omitempty"`
	Sheets      []Sheet `firestore:"sheets" json:"sheets,omitempty"`
}

type SheetConnector struct {
	Db             *FireDb
	CollectionName string
}

func (c *SheetConnector) SearchByGroupKeys(groupKeys []string) ([]Sheet, error) {
	if len(groupKeys) == 0 {
		return []Sheet{}, nil
	}
	docs, err := c.Db.SearchByQuery(c.CollectionName, GroupKeysQuery(groupKeys))
	if err != nil {
		return nil, err
	}
	groups, err := ParseSheetGroupCollection(docs)
	if err != nil {
		return nil, err
	}
	return FlattenSheets(groups), nil
}

func (c *SheetConnector) GetAllTitles() ([]SheetGroup, error) {
	docs, err := c.Db.SearchByQuery(c.CollectionName, GetTitles())
	if err != nil {
		return nil, err
	}
	return ParseSheetGroupCollection(docs)
}

func (c *SheetConnector) AddSheet(sheet *Sheet) error {
	return c.Db.UpdateDocumentWithGroupKey(c.CollectionName, GroupKey(sheet.Band, sheet.Title), AddSheetFn(sheet))
}

func (c *SheetConnector) UpsertSheet(sheet *Sheet) error {
	return c.Db.UpdateDocumentWithGroupKey(c.CollectionName, GroupKey(sheet.Band, sheet.Title), UpsertSheetFn(sheet))
}

func (c *SheetConnector) RemoveSheet(sheet *Sheet) error {
	return c.Db.UpdateDocumentWithGroupKey(c.CollectionName, GroupKey(sheet.Band, sheet.Title), RemoveSheetFn(sheet))
}

func (sg *SheetGroup) AddSheetIfMissing(sheet *Sheet) {
	if sg == nil {
		return
	}
	for i := range sg.Sheets {
		if sg.Sheets[i].Id == sheet.Id {
			sg.Sheets[i] = *sheet
			return
		}
	}
	sg.Sheets = append(sg.Sheets, *sheet)
}

func GetTitles() DocQuery {
	return func(cr *firestore.CollectionRef) firestore.Query {
		return cr.SelectPaths(firestore.FieldPath{firestore.DocumentID}, firestore.FieldPath{"title"}, firestore.FieldPath{"band"}, firestore.FieldPath{"group_key"})
	}
}

func AddSheetFn(sheet *Sheet) DocUpdateFn {
	return func(transaction *firestore.Transaction, ref *firestore.DocumentRef) error {
		err := transaction.Set(ref, map[string]interface{}{"sheets": firestore.ArrayUnion(sheet),
			"title": sheet.Title, "band": sheet.Band, "group_key": GroupKey(sheet.Band, sheet.Title), "last_updated": time.Now().UTC().String()}, firestore.MergeAll)
		return err
	}
}

func UpsertSheetFn(sheet *Sheet) DocUpdateFn {
	return func(transaction *firestore.Transaction, ref *firestore.DocumentRef) error {
		var sg SheetGroup
		if doc, err := transaction.Get(ref); err == nil {
			if err := doc.DataTo(&sg); err != nil {
				return err
			}
		}
		sg.AddSheetIfMissing(sheet)
		sg.Title = sheet.Title
		sg.Band = sheet.Band
		sg.GroupKey = GroupKey(sheet.Band, sheet.Title)
		sg.LastUpdated = time.Now().UTC().String()
		return transaction.Set(ref, sg)
	}
}

func RemoveSheetFn(sheet *Sheet) DocUpdateFn {
	return func(transaction *firestore.Transaction, ref *firestore.DocumentRef) error {
		var sg SheetGroup
		doc, err := transaction.Get(ref)
		if err != nil {
			return err
		}
		if err := doc.DataTo(&sg); err != nil {
			return err
		}

		filtered := sg.Sheets[:0]
		for _, s := range sg.Sheets {
			if s.Id != sheet.Id {
				filtered = append(filtered, s)
			}
		}

		if len(filtered) == 0 {
			return transaction.Delete(ref)
		}

		sg.Sheets = filtered
		sg.LastUpdated = time.Now().UTC().String()
		return transaction.Set(ref, sg)
	}
}

func (sg *SheetGroup) GetSheets() []Sheet {
	baseURL := strings.TrimRight(os.Getenv(sheetStorageBaseURLEnv), "/")
	bucketName := strings.TrimSpace(os.Getenv(sheetStorageBucketEnv))
	for i := range sg.Sheets {
		sg.Sheets[i].Title = sg.Title
		sg.Sheets[i].Band = sg.Band
		if sg.Sheets[i].Id != "" && baseURL != "" && bucketName != "" {
			sg.Sheets[i].Url = baseURL + "/storage/v1/b/" + bucketName + "/o/" + sg.Sheets[i].Id
		}
	}
	return sg.Sheets
}

func ToSheet(rawPath string, metadata map[string]string) (*Sheet, error) {
	fileName := path.Base(rawPath)
	title := path.Base(path.Dir(rawPath))
	band := path.Base(path.Dir(path.Dir(rawPath)))
	if title == "." || title == "/" {
		title = ""
	}
	if band == "." || band == "/" {
		band = ""
	}
	instrumentName := metadata["instrument_name_en"]
	instrumentNamePol := metadata["instrument_name_pol"]
	voice := metadata["voice"]
	key := metadata["key"]
	id := rawPath

	instrument := Instrument{
		Name: instrumentName, NamePol: instrumentNamePol, Key: key, Voice: voice}
	return &Sheet{Instrument: instrument, Id: id, FileName: fileName, Title: title, Band: band}, nil
}

func NewSheetGroup(title string, sheet Sheet) *SheetGroup {
	return &SheetGroup{
		Title:       title,
		Band:        sheet.Band,
		GroupKey:    GroupKey(sheet.Band, title),
		LastUpdated: time.Now().UTC().String(),
		Sheets:      []Sheet{sheet},
	}
}

func GroupKey(band, title string) string {
	normalizedBand := normalizeGroupPart(band)
	normalizedTitle := normalizeGroupPart(title)
	if normalizedBand == "" {
		return normalizedTitle
	}
	if normalizedTitle == "" {
		return normalizedBand
	}
	return normalizedBand + "__" + normalizedTitle
}

func normalizeGroupPart(value string) string {
	lowered := strings.ToLower(strings.TrimSpace(value))
	if lowered == "" {
		return ""
	}
	collapsed := strings.Join(strings.Fields(lowered), " ")
	slug := groupKeySeparatorRegex.ReplaceAllString(collapsed, "-")
	return strings.Trim(slug, "-")
}

func ParseSheetGroup(doc *firestore.DocumentSnapshot) (*SheetGroup, error) {
	var sg SheetGroup
	if err := doc.DataTo(&sg); err != nil {
		return nil, err
	}
	sg.ID = doc.Ref.ID
	return &sg, nil
}

func ParseSheetGroupCollection(docs []*firestore.DocumentSnapshot) ([]SheetGroup, error) {
	var groups []SheetGroup
	for _, doc := range docs {
		sg, err := ParseSheetGroup(doc)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *sg)
	}
	return groups, nil
}

func FlattenSheets(sheetGroups []SheetGroup) []Sheet {
	result := make([]Sheet, 0)
	for _, group := range sheetGroups {
		result = append(result, group.GetSheets()...)
	}
	return result
}
