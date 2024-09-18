package notiontest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	mrand "math/rand"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/jarcoal/httpmock"
	"golang.org/x/oauth2"

	"go.f110.dev/notion-api/v3"
)

const (
	secretSeed = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type versionHandler interface {
	ListAllUsers(req *http.Request) (*http.Response, error)
	GetUser(req *http.Request) (*http.Response, error)
	CreateDatabase(req *http.Request) (*http.Response, error)
	GetDatabase(req *http.Request) (*http.Response, error)
	UpdateDatabase(req *http.Request) (*http.Response, error)
}

type Mock struct {
	users     []*notion.User
	databases []*notion.Database
	tokens    map[string]*notion.User
}

// NewMock returns the mock object for Notion API.
func NewMock() *Mock {
	return &Mock{tokens: make(map[string]*notion.User)}
}

func (n *Mock) RegisterMock(mock *httpmock.MockTransport) {
	n.registerUsers(mock)
	n.registerDatabases(mock)
}

// AuthenticatedClient returns a http client of the authenticated bot user.
// If you need raw http.RoundTripper with API mock, then you will need to use RegisterMock.
func (n *Mock) AuthenticatedClient(botName string) *http.Client {
	tr := httpmock.NewMockTransport()
	n.RegisterMock(tr)

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: n.GenerateBotToken(botName)})
	tc := oauth2.NewClient(context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: tr}), ts)
	return tc
}

// User adds new user for human.
func (n *Mock) User(name string) *Mock {
	n.users = append(n.users, &notion.User{Meta: &notion.Meta{ID: newID(), Object: "user"}, Type: notion.UserTypePerson, Name: name})
	return n
}

// BotUser adds new user for a machine.
func (n *Mock) BotUser(name string) *Mock {
	n.users = append(n.users, &notion.User{Meta: &notion.Meta{ID: newID(), Object: "user"}, Type: notion.UserTypeBot, Name: name})
	return n
}

// Database adds a database.
func (n *Mock) Database(db *notion.Database) *Mock {
	if db.ID == "" {
		if db.Meta == nil {
			db.Meta = &notion.Meta{}
		}
		db.ID = newID()
	}
	n.databases = append(n.databases, db)
	return n
}

// NewDatabase creates and adds a database.
func (n *Mock) NewDatabase(title string) *notion.Database {
	db := NewDatabase(title)
	n.Database(db)
	return db
}

// FindUser returns notion.User with a name.
func (n *Mock) FindUser(name string) *notion.User {
	for _, user := range n.users {
		if user.Name == name {
			return user
		}
	}
	return nil
}

// FindDatabase returns the list of notion.Database that has specified title.
func (n *Mock) FindDatabase(title string) []*notion.Database {
	var dbs []*notion.Database
	for _, db := range n.databases {
		if len(db.Title) == 0 {
			continue
		}
		if db.Title[0].Text.Content == title {
			dbs = append(dbs, db)
		}
	}
	return dbs
}

// GenerateBotToken generates a new token for the bot with a name
func (n *Mock) GenerateBotToken(botName string) string {
	var bot *notion.User
	for _, v := range n.users {
		if v.Name == botName {
			bot = v
			break
		}
	}
	if bot == nil {
		return ""
	}

	buf := make([]byte, 43)
	for i := 0; i < len(buf); i++ {
		buf[i] = secretSeed[mrand.Intn(len(secretSeed))]
	}
	token := "secret_" + string(buf)
	n.tokens[token] = bot
	return token
}

func (n *Mock) registerUsers(mock *httpmock.MockTransport) {
	// List all users
	n.registerResponderForAuthorizedRequest(mock,
		http.MethodGet,
		regexp.MustCompile(`/v1/users$`),
		func(req *http.Request, handler versionHandler) (*http.Response, error) {
			return handler.ListAllUsers(req)
		},
	)
	// Get the user
	n.registerResponderForAuthorizedRequest(mock,
		http.MethodGet,
		regexp.MustCompile(`/v1/users/[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$`),
		func(req *http.Request, handler versionHandler) (*http.Response, error) {
			return handler.GetUser(req)
		},
	)
	// Get me
	n.registerResponderForAuthorizedRequest(mock,
		http.MethodGet,
		regexp.MustCompile(`/v1/users/me$`),
		func(req *http.Request, handler versionHandler) (*http.Response, error) {
			token := n.getToken(req)
			user := n.tokens[token]
			return httpmock.NewJsonResponse(http.StatusOK, user)
		},
	)
}

func (n *Mock) registerDatabases(mock *httpmock.MockTransport) {
	// Create a database
	n.registerResponderForAuthorizedRequest(mock,
		http.MethodPost,
		regexp.MustCompile(`/v1/databases$`),
		func(req *http.Request, handler versionHandler) (*http.Response, error) {
			return handler.CreateDatabase(req)
		},
	)
	// Get a database
	n.registerResponderForAuthorizedRequest(mock,
		http.MethodGet,
		regexp.MustCompile(`/v1/databases/[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$`),
		func(req *http.Request, handler versionHandler) (*http.Response, error) {
			return handler.GetDatabase(req)
		},
	)
	// Update a database
	n.registerResponderForAuthorizedRequest(mock,
		http.MethodPatch,
		regexp.MustCompile(`/v1/databases/[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$`),
		func(req *http.Request, handler versionHandler) (*http.Response, error) {
			return handler.UpdateDatabase(req)
		},
	)
	// Query a database
	// TODO: This isn't implemented yet.
}

func (n *Mock) registerResponderForAuthorizedRequest(mock *httpmock.MockTransport, method string, urlRegexp *regexp.Regexp, responder func(req *http.Request, h versionHandler) (*http.Response, error)) {
	mock.RegisterRegexpResponder(method, urlRegexp, func(req *http.Request) (*http.Response, error) {
		if !n.authorizeRequest(req) {
			return n.unauthorizedError()
		}
		handler := n.getHandler(req.Header)
		if handler == nil {
			return n.missionVersionError()
		}

		return responder(req, handler)
	})
}

func (n *Mock) authorizeRequest(req *http.Request) bool {
	token := n.getToken(req)
	if _, ok := n.tokens[token]; !ok {
		return false
	}
	return true
}

func (n *Mock) getToken(req *http.Request) string {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	idx := strings.Index(authHeader, "Bearer ")
	return authHeader[idx+len("Bearer "):]
}

func (n *Mock) getHandler(header http.Header) versionHandler {
	switch header.Get("Notion-Version") {
	case "2022-06-28":
		return &version220628{m: n}
	default:
		return nil
	}
}

func (n *Mock) missionVersionError() (*http.Response, error) {
	e := &notion.Error{
		Meta: &notion.Meta{
			Object: "object",
		},
		Status:  400,
		Code:    "missing_version",
		Message: "Notion-Version header failed validation: Notion-Version header should be defined, instead was `undefined`.",
	}
	return httpmock.NewJsonResponse(e.Status, e)
}

func (n *Mock) unauthorizedError() (*http.Response, error) {
	e := &notion.Error{
		Meta: &notion.Meta{
			Object: "object",
		},
		Status:  401,
		Code:    "unauthorized",
		Message: "API token is invalid.",
	}
	return httpmock.NewJsonResponse(e.Status, e)
}

type version220628 struct {
	m *Mock
}

var _ versionHandler = (*version220628)(nil)

func (h *version220628) ListAllUsers(req *http.Request) (*http.Response, error) {
	results, hasMore, nextCursor := sliceByPagination(req, h.m.users)
	res := &notion.UserList{
		ListMeta: &notion.ListMeta{
			Object:     "list",
			HasMore:    hasMore,
			NextCursor: nextCursor,
		},
		Results: results,
	}
	return httpmock.NewJsonResponse(http.StatusOK, res)
}

func (h *version220628) GetUser(req *http.Request) (*http.Response, error) {
	_, userID := path.Split(req.URL.Path)
	for _, v := range h.m.users {
		if v.GetID() == userID {
			return httpmock.NewJsonResponse(http.StatusOK, v)
		}
	}

	e := newErrorResponse(400, "object_not_found", fmt.Sprintf("Could not find user with ID: %s.", userID))
	return httpmock.NewJsonResponse(e.Status, e)
}

func (h *version220628) CreateDatabase(req *http.Request) (*http.Response, error) {
	d := json.NewDecoder(req.Body)
	var db *notion.Database
	if err := d.Decode(&db); err != nil {
		return httpmock.NewJsonResponse(http.StatusBadRequest, nil)
	}

	if db.Meta == nil {
		db.Meta = &notion.Meta{}
	}
	db.ID = newID()
	h.m.Database(db)
	return httpmock.NewJsonResponse(http.StatusOK, db)
}

func (h *version220628) GetDatabase(req *http.Request) (*http.Response, error) {
	_, databaseID := path.Split(req.URL.Path)
	for _, v := range h.m.databases {
		if v.GetID() == databaseID {
			return httpmock.NewJsonResponse(http.StatusOK, v)
		}
	}

	e := newErrorResponse(404, "object_not_found", fmt.Sprintf("Could not find database with ID: %s. Make sure the relevant pages and databases are shared with your integration.", databaseID))
	return httpmock.NewJsonResponse(e.Status, e)
}

func (h *version220628) UpdateDatabase(req *http.Request) (*http.Response, error) {
	_, databaseID := path.Split(req.URL.Path)
	var db *notion.Database
	for _, v := range h.m.databases {
		if v.GetID() == databaseID {
			db = v
			break
		}
	}
	if db == nil {
		e := newErrorResponse(404, "object_not_found", fmt.Sprintf("Could not find database with ID: %s. Make sure the relevant pages and databases are shared with your integration.", databaseID))
		return httpmock.NewJsonResponse(e.Status, e)
	}

	var updatedDB *notion.Database
	if err := json.NewDecoder(req.Body).Decode(&updatedDB); err != nil {
		return httpmock.NewJsonResponse(http.StatusBadRequest, nil)
	}
	updatedDB.Meta = db.Meta
	for i, v := range h.m.databases {
		if v.GetID() == db.GetID() {
			h.m.databases = h.m.databases[:i]
			h.m.databases = append(h.m.databases, updatedDB)
			h.m.databases = append(h.m.databases, h.m.databases[i+1:]...)
			break
		}
	}

	return httpmock.NewJsonResponse(http.StatusOK, updatedDB)
}

type abstractObject interface {
	GetID() string
}

func sliceByPagination[T abstractObject](req *http.Request, s []T) ([]T, bool, string) {
	startCursor, pageSize := getPagination(req)
	startIdx := 0
	if startCursor != "" {
		for i, v := range s {
			if v.GetID() == startCursor {
				startIdx = i
				break
			}
		}
	}
	endIdx := startIdx + pageSize
	hasMore := true
	if endIdx >= len(s) {
		endIdx = len(s)
		hasMore = false
	}
	var nextCursor string
	if hasMore {
		nextCursor = s[endIdx].GetID()
	}

	return s[startIdx:endIdx], hasMore, nextCursor
}

func getPagination(req *http.Request) (startCursor string, pageSize int) {
	// Default page size is 100
	pageSize = 100
	size := req.URL.Query().Get("page_size")
	if size != "" {
		pageSize, _ = strconv.Atoi(size)
	}
	return req.URL.Query().Get("start_cursor"), pageSize
}

func newID() string {
	buf := make([]byte, 16)
	io.ReadFull(rand.Reader, buf)

	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80

	st := make([]byte, 36)
	hex.Encode(st, buf[:4])
	st[8] = '-'
	hex.Encode(st[9:13], buf[4:6])
	st[13] = '-'
	hex.Encode(st[14:18], buf[6:8])
	st[18] = '-'
	hex.Encode(st[19:23], buf[8:10])
	st[23] = '-'
	hex.Encode(st[24:], buf[10:])
	return string(st)
}

func newErrorResponse(status int, code, msg string) *notion.Error {
	return &notion.Error{
		Meta: &notion.Meta{
			Object: "object",
		},
		Status:  status,
		Code:    code,
		Message: msg,
	}
}
