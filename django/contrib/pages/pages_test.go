package pages_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/Nigel2392/django/contrib/pages"
	"github.com/Nigel2392/django/contrib/pages/models"
)

var sqlDB *sql.DB

func getEnv(key, def string) string {
	var val = def
	if v, ok := os.LookupEnv(key); ok {
		val = v
	}
	return val
}

func init() {

	var (
		dbEngine = getEnv("DB_ENGINE", "sqlite3")
		dbURL    = getEnv("DB_URL", "file::memory:?cache=shared")
	)

	var err error
	sqlDB, err = sql.Open(dbEngine, dbURL)
	if err != nil {
		panic(err)
	}

	// Create test_pages table
	if err := sqlDB.Ping(); err != nil {
		panic(err)
	}

	if _, err := sqlDB.Exec(testPageCREATE_TABLE); err != nil {
		panic(err)
	}
}

var _ pages.Page = &TestPage{}

type TestPage struct {
	Ref         *models.PageNode
	Identifier  int64
	Description string
}

func nodesEqual(a, b *models.PageNode) bool {
	pageARef := *a
	pageBRef := *b
	pageARef.CreatedAt = pageBRef.CreatedAt
	pageARef.UpdatedAt = pageBRef.UpdatedAt
	return pageARef == pageBRef

}

func (t *TestPage) ID() int64 {
	return int64(t.Identifier)
}

func (t *TestPage) Reference() *models.PageNode {
	return t.Ref
}

func (t *TestPage) Save(ctx context.Context) error {
	return nil
}

type DBTestPage struct {
	TestPage
}

func (t *DBTestPage) Save(ctx context.Context) error {
	if t.Identifier == 0 {
		result, err := sqlDB.ExecContext(ctx, testPageINSERT, t.Description)
		if err != nil {
			return err
		}
		t.Identifier, err = result.LastInsertId()
		return err
	}
	_, err := sqlDB.ExecContext(ctx, testPageUPDATE, t.Description, t.Identifier)
	return err
}

const (
	testPageINSERT       = `INSERT INTO test_pages (title) VALUES (?)`
	testPageUPDATE       = `UPDATE test_pages SET title = ? WHERE id = ?`
	testPageByID         = `SELECT id, title FROM test_pages WHERE id = ?`
	testPageCREATE_TABLE = `CREATE TABLE IF NOT EXISTS test_pages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT
	)`
)

func TestContentType(t *testing.T) {
	var cType = pages.NewContentType(&TestPage{})

	t.Run("NewPage", func(t *testing.T) {
		var newInstance = cType.New()
		var _, ok = newInstance.(*TestPage)
		if !ok {
			t.Errorf("expected *TestPage, got %T", newInstance)
		}
	})

	t.Run("TypeName", func(t *testing.T) {
		if cType.Model() != "TestPage" {
			t.Errorf("expected TestPage as Model, got %s", cType.Model())
		}
	})

}

func TestPageRegistry(t *testing.T) {
	t.Run("Definition", func(t *testing.T) {

		var definition = &pages.PageDefinition{
			PageObject: &TestPage{
				Ref: &models.PageNode{
					ID:     1,
					Title:  "Title",
					Path:   "001",
					PageID: 69,
					Depth:  0,
				},
				Identifier:  69,
				Description: "Description",
			},
			GetForID: func(ctx context.Context, node models.PageNode, id int64) (pages.SaveablePage, error) {
				return &TestPage{
					Ref:         &node,
					Identifier:  id,
					Description: fmt.Sprintf("Description for %s (%d)", node.Title, id),
				}, nil
			},
		}

		pages.Register(definition)
		var cType = definition.ContentType()

		t.Run("SpecificInstance", func(t *testing.T) {
			var node = models.PageNode{
				Title:       "Test Page",
				PageID:      69,
				ContentType: cType.TypeName(),
			}
			var instance, err = pages.Specific(context.Background(), node)
			if err != nil {
				t.Error(err)
			}
			var testPage, ok = instance.(*TestPage)
			if !ok {
				t.Errorf("expected *TestPage, got %T", instance)
			}
			var descString = fmt.Sprintf("Description for %s (%d)", node.Title, node.PageID)
			if testPage.Description != descString {
				t.Errorf("expected %q, got %q", descString, testPage.Description)
			}
		})

		t.Run("ForObject", func(t *testing.T) {
			var testDef = pages.DefinitionForObject(&TestPage{})
			if testDef == nil {
				t.Error("expected definition, got nil")
				return
			}

			if definition.PageObject != testDef.PageObject {
				t.Errorf("expected %+v, got %+v", definition.PageObject, testDef.PageObject)
			}

			var page = testDef.PageObject.(*TestPage)
			if page.Description != "Description" {
				t.Errorf("expected Description, got %s", page.Description)
			}

			if page.Identifier != 69 {
				t.Errorf("expected Identifier, got %d", page.Identifier)
			}

			if page.Ref.Title != "Title" {
				t.Errorf("expected Title, got %s", page.Ref.Title)
			}
		})
	})
}

func TestPageNode(t *testing.T) {
	var (
		rootNode = models.PageNode{
			Title: "Root",
		}
		childNode = models.PageNode{
			Title: "Child",
		}
		childSiblingNode = models.PageNode{
			Title: "ChildSibling",
		}
		subChildNode = models.PageNode{
			Title: "SubChild",
		}
		childSiblingSubChildNode = models.PageNode{
			Title: "ChildSiblingSubChild",
		}
		queryCtx = context.Background()
		querier  = pages.QuerySet(sqlDB)
	)

	if err := pages.CreateTable(sqlDB); err != nil {
		t.Error(err)
		return
	}

	t.Run("Root", func(t *testing.T) {
		var err = pages.CreateRootNode(querier, queryCtx, &rootNode)
		if err != nil {
			t.Error(err)
			return
		}

		if rootNode.ID != 1 {
			t.Errorf("expected ID 1, got %d", rootNode.ID)
		}

		if rootNode.Path != "001" {
			t.Errorf("expected Path 1, got %s", rootNode.Path)
		}

		if rootNode.Depth != 0 {
			t.Errorf("expected Depth 0, got %d", rootNode.Depth)
		}

		if rootNode.Numchild != 0 {
			t.Errorf("expected Numchild 0, got %d", rootNode.Numchild)
		}

		if rootNode.StatusFlags != 0 {
			t.Errorf("expected StatusFlagPublished, got %d", rootNode.StatusFlags)
		}

		if rootNode.PageID != 0 {
			t.Errorf("expected PageID 0, got %d", rootNode.PageID)
		}

		if rootNode.ContentType != "" {
			t.Errorf("expected ContentType empty, got %s", rootNode.ContentType)
		}

		t.Run("AddChild", func(t *testing.T) {
			var err = pages.CreateChildNode(querier, queryCtx, &rootNode, &childNode)
			if err != nil {
				t.Error(err)
				return
			}

			if childNode.ID != 2 {
				t.Errorf("expected ID 2, got %d", childNode.ID)
			}

			if childNode.Path != "001001" {
				t.Errorf("expected Path 001001, got %s", childNode.Path)
			}

			if childNode.Depth != 1 {
				t.Errorf("expected Depth 1, got %d", childNode.Depth)
			}

			if childNode.Numchild != 0 {
				t.Errorf("expected Numchild 0, got %d", childNode.Numchild)
			}

			if childNode.StatusFlags != 0 {
				t.Errorf("expected StatusFlagPublished, got %d", childNode.StatusFlags)
			}

			if childNode.PageID != 0 {
				t.Errorf("expected PageID 0, got %d", childNode.PageID)
			}

			if childNode.ContentType != "" {
				t.Errorf("expected ContentType empty, got %s", childNode.ContentType)
			}

			if rootNode.Numchild != 1 {
				t.Errorf("expected Numchild 1, got %d", rootNode.Numchild)
			}

			t.Run("GetChildren", func(t *testing.T) {
				var children, err = querier.GetChildNodes(queryCtx, rootNode.Path, rootNode.Depth)
				if err != nil {
					t.Error(err)
					return
				}

				if len(children) != 1 {
					t.Errorf("expected 1 child, got %d", len(children))
					return
				}

				if !nodesEqual(&children[0], &childNode) {
					t.Errorf("expected %+v, got %+v", childNode, children[0])
					return
				}
			})

			t.Run("AddSubChild", func(t *testing.T) {
				var err = pages.CreateChildNode(querier, queryCtx, &childNode, &subChildNode)
				if err != nil {
					t.Error(err)
					return
				}

				if subChildNode.ID != 3 {
					t.Errorf("expected ID 3, got %d", subChildNode.ID)
				}

				if subChildNode.Path != "001001001" {
					t.Errorf("expected Path 001001001, got %s", subChildNode.Path)
				}

				if subChildNode.Depth != 2 {
					t.Errorf("expected Depth 2, got %d", subChildNode.Depth)
				}

				if subChildNode.Numchild != 0 {
					t.Errorf("expected Numchild 0, got %d", subChildNode.Numchild)
				}

				if childNode.Numchild != 1 {
					t.Errorf("expected Numchild 1, got %d", childNode.Numchild)
				}

				t.Run("GetAncestors", func(t *testing.T) {
					var ancestors, err = pages.AncestorNodes(querier, queryCtx, subChildNode.Path, int(subChildNode.Depth)+1)
					if err != nil {
						t.Error(err)
						return
					}

					if len(ancestors) != 2 {
						t.Errorf("expected 2 ancestors, got %d", len(ancestors))
						return
					}

					if !nodesEqual(&ancestors[0], &rootNode) {
						t.Errorf("expected %+v, got %+v", rootNode, ancestors[0])
						return
					}

					if !nodesEqual(&ancestors[1], &childNode) {
						t.Errorf("expected %+v, got %+v", childNode, ancestors[1])
						return
					}
				})

				t.Run("GetRootDescendants", func(t *testing.T) {
					var descendants, err = querier.GetDescendants(queryCtx, rootNode.Path, 0)
					if err != nil {
						t.Error(err)
						return
					}

					if len(descendants) != 2 {
						t.Errorf("expected 2 descendants, got %d", len(descendants))
						return
					}

					if !nodesEqual(&descendants[0], &childNode) {
						t.Errorf("expected %+v, got %+v", childNode, descendants[0])
						return
					}

					if !nodesEqual(&descendants[1], &subChildNode) {
						t.Errorf("expected %+v, got %+v", subChildNode, descendants[1])
						return
					}
				})

				t.Run("ParentNode", func(t *testing.T) {
					var parent, err = pages.ParentNode(querier, queryCtx, subChildNode.Path, int(subChildNode.Depth))
					if err != nil {
						t.Error(err)
						return
					}

					if !nodesEqual(&parent, &childNode) {
						t.Errorf("expected %+v, got %+v", childNode, parent)
					}
				})

				t.Run("DeleteNode", func(t *testing.T) {
					var err = pages.DeleteNode(querier, queryCtx, subChildNode.ID, subChildNode.Path, subChildNode.Depth)
					if err != nil {
						t.Error(err)
						return
					}

					descendants, err := querier.GetDescendants(queryCtx, rootNode.Path, 0)
					if err != nil {
						t.Error(err)
						return
					}

					if len(descendants) != 1 {
						t.Errorf("expected 1 descendant, got %d", len(descendants))
						return
					}

					childNode, err = querier.GetNodeByID(queryCtx, childNode.ID)
					if err != nil {
						t.Error(err)
						return
					}

					if !nodesEqual(&descendants[0], &childNode) {
						t.Errorf("expected %+v, got %+v", childNode, descendants[0])
						return
					}

					if childNode.Numchild != 0 {
						t.Errorf("expected Numchild 0, got %d", childNode.Numchild)
						return
					}
				})
			})
		})

		t.Run("AddSibling", func(t *testing.T) {
			var err = pages.CreateChildNode(querier, queryCtx, &rootNode, &childSiblingNode)
			if err != nil {
				t.Error(err)
				return
			}

			if childSiblingNode.ID != 4 {
				t.Errorf("expected ID 3, got %d", childSiblingNode.ID)
			}

			if childSiblingNode.Path != "001002" {
				t.Errorf("expected Path 001002, got %s", childSiblingNode.Path)
			}

			if childSiblingNode.Depth != 1 {
				t.Errorf("expected Depth 1, got %d", childSiblingNode.Depth)
			}

			if childSiblingNode.Numchild != 0 {
				t.Errorf("expected Numchild 0, got %d", childSiblingNode.Numchild)
			}

			if childSiblingNode.StatusFlags != 0 {
				t.Errorf("expected StatusFlagPublished, got %d", childSiblingNode.StatusFlags)
			}

			if childSiblingNode.PageID != 0 {
				t.Errorf("expected PageID 0, got %d", childSiblingNode.PageID)
			}

			if childSiblingNode.ContentType != "" {
				t.Errorf("expected ContentType empty, got %s", childSiblingNode.ContentType)
			}

			if rootNode.Numchild != 2 {
				t.Errorf("expected Numchild 2, got %d", rootNode.Numchild)
			}

			t.Run("GetChildren", func(t *testing.T) {
				var children, err = querier.GetChildNodes(queryCtx, rootNode.Path, rootNode.Depth)
				if err != nil {
					t.Error(err)
					return
				}

				if len(children) != 2 {
					t.Errorf("expected 2 children, got %d", len(children))
					return
				}

				if !nodesEqual(&children[1], &childSiblingNode) {
					t.Errorf("expected %+v, got %+v", childSiblingNode, children[1])
					return
				}
			})

			t.Run("AddSubChild", func(t *testing.T) {
				var err = pages.CreateChildNode(querier, queryCtx, &childSiblingNode, &childSiblingSubChildNode)
				if err != nil {
					t.Error(err)
					return
				}

				if childSiblingSubChildNode.ID != 5 {
					t.Errorf("expected ID 5, got %d", childSiblingSubChildNode.ID)
				}

				if childSiblingSubChildNode.Path != "001002001" {
					t.Errorf("expected Path 001002001, got %s", childSiblingSubChildNode.Path)
				}

				if childSiblingSubChildNode.Depth != 2 {
					t.Errorf("expected Depth 2, got %d", childSiblingSubChildNode.Depth)
				}

				if childSiblingSubChildNode.Numchild != 0 {
					t.Errorf("expected Numchild 0, got %d", childSiblingSubChildNode.Numchild)
				}

				if childSiblingNode.Numchild != 1 {
					t.Errorf("expected Numchild 1, got %d", childSiblingNode.Numchild)
				}

				t.Run("GetAncestors", func(t *testing.T) {
					var ancestors, err = pages.AncestorNodes(querier, queryCtx, childSiblingSubChildNode.Path, int(childSiblingSubChildNode.Depth)+1)
					if err != nil {
						t.Error(err)
						return
					}

					if len(ancestors) != 2 {
						t.Errorf("expected 2 ancestors, got %d", len(ancestors))
						return
					}

					if !nodesEqual(&ancestors[0], &rootNode) {
						t.Errorf("expected %+v, got %+v", rootNode, ancestors[0])
						return
					}

					if !nodesEqual(&ancestors[1], &childSiblingNode) {
						t.Errorf("expected %+v, got %+v", childSiblingNode, ancestors[1])
						return
					}
				})

				t.Run("GetRootDescendants", func(t *testing.T) {
					var descendants, err = querier.GetDescendants(queryCtx, rootNode.Path, 0)
					if err != nil {
						t.Error(err)
						return
					}

					if len(descendants) != 3 {
						t.Errorf("expected 3 descendants, got %d", len(descendants))
						return
					}

					if !nodesEqual(&descendants[1], &childSiblingNode) {
						t.Errorf("expected %+v, got %+v", childSiblingNode, descendants[1])
						return
					}

					if !nodesEqual(&descendants[2], &childSiblingSubChildNode) {
						t.Errorf("expected %+v, got %+v", childSiblingSubChildNode, descendants[2])
						return
					}
				})
			})
		})
	})

	pages.Register(&pages.PageDefinition{
		PageObject: &DBTestPage{},
		GetForID: func(ctx context.Context, ref models.PageNode, id int64) (pages.SaveablePage, error) {
			var page = &DBTestPage{}
			page.Ref = &ref
			var row = sqlDB.QueryRowContext(ctx, testPageByID, id)
			if err := row.Scan(&page.Identifier, &page.Description); err != nil {
				return nil, err
			}
			return page, nil
		},
	})

	var createAndBind = func(t *testing.T, node *models.PageNode, page *DBTestPage) *DBTestPage {
		page.Ref = node
		if err := page.Save(queryCtx); err != nil {
			t.Error(err)
		}
		node.PageID = page.ID()
		node.ContentType = pages.DefinitionForObject(page).ContentType().TypeName()

		if err := querier.UpdateNode(queryCtx, node.Title, node.Path, node.Depth, node.Numchild, int64(node.StatusFlags), node.PageID, node.ContentType, node.ID); err != nil {
			t.Error(err)
		}

		return page
	}

	var (
		dbTestPage_Root = createAndBind(t, &rootNode, &DBTestPage{
			TestPage: TestPage{
				Description: "Root Description",
			},
		})
		dbTestPage_Child = createAndBind(t, &childNode, &DBTestPage{
			TestPage: TestPage{
				Description: "Child Description",
			},
		})
		dbTestPage_ChildSibling = createAndBind(t, &childSiblingNode, &DBTestPage{
			TestPage: TestPage{
				Description: "ChildSibling Description",
			},
		})
		dbTestPage_SubChild = createAndBind(t, &subChildNode, &DBTestPage{
			TestPage: TestPage{
				Description: "SubChild Description",
			},
		})
		dbTestPage_ChildSiblingSubChild = createAndBind(t, &childSiblingSubChildNode, &DBTestPage{
			TestPage: TestPage{
				Description: "ChildSiblingSubChild Description",
			},
		})
	)

	var pageList = []pages.Page{
		dbTestPage_Root,
		dbTestPage_Child,
		dbTestPage_ChildSibling,
		dbTestPage_SubChild,
		dbTestPage_ChildSiblingSubChild,
	}

	for _, page := range pageList {
		t.Run(fmt.Sprintf("Specific_Page_%s", page.Reference().Title), func(t *testing.T) {
			var instance, err = pages.Specific(queryCtx, *page.Reference())
			if err != nil {
				t.Error(err)
				return
			}

			var dbPage, ok = instance.(*DBTestPage)
			if !ok {
				t.Errorf("expected *DBTestPage, got %T", instance)
				return
			}

			if dbPage.Description != page.(*DBTestPage).Description {
				t.Errorf("expected %s, got %s", page.(*DBTestPage).Description, dbPage.Description)
			}
		})
	}
}
