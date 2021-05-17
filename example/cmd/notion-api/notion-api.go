package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"go.f110.dev/notion-api"
	"golang.org/x/oauth2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "notion-api action [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get-user":
		if err := getUser(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "list-users":
		if err := listUsers(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "list-databases":
		if err := listDatabases(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "get-database":
		if err := getDatabase(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "get-pages":
		if err := getPages(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "get-blocks":
		if err := getBlocks(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "create-page":
		if err := createPage(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "update-properties":
		if err := updateProperties(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "append-blocks":
		if err := appendBlocks(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "%s is not action\n", os.Args[1])
		os.Exit(1)
	}
}

func getUser(args []string) error {
	userID := ""
	token := ""
	fs := flag.NewFlagSet("get-user", flag.ContinueOnError)
	fs.StringVar(&userID, "user-id", "", "User identifier")
	fs.StringVar(&token, "token", "", "API Token")
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := newClient(token)
	if err != nil {
		return err
	}

	user, err := client.GetUser(context.Background(), userID)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", user)

	return nil
}

func listUsers(args []string) error {
	token := ""
	fs := flag.NewFlagSet("list-users", flag.ContinueOnError)
	fs.StringVar(&token, "token", "", "API Token")
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := newClient(token)
	if err != nil {
		return err
	}
	users, err := client.ListAllUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		fmt.Printf("ID: %s %+v\n", user.ID, user)
	}

	return nil
}

func listDatabases(args []string) error {
	token := ""
	fs := flag.NewFlagSet("list-databases", flag.ContinueOnError)
	fs.StringVar(&token, "token", "", "API Token")
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := newClient(token)
	if err != nil {
		return err
	}
	databases, err := client.ListDatabases(context.Background())
	if err != nil {
		return err
	}
	for _, database := range databases {
		fmt.Printf("ID: %s %+v\n", database.ID, database)
	}

	return nil
}

func getDatabase(args []string) error {
	databaseID := ""
	token := ""
	fs := flag.NewFlagSet("get-database", flag.ContinueOnError)
	fs.StringVar(&databaseID, "database-id", "", "Database identifier")
	fs.StringVar(&token, "token", "", "API Token")
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := newClient(token)
	if err != nil {
		return err
	}

	database, err := client.GetDatabase(context.Background(), databaseID)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", database)

	return nil
}

func getPages(args []string) error {
	databaseID := ""
	token := ""
	fs := flag.NewFlagSet("get-pages", flag.ContinueOnError)
	fs.StringVar(&databaseID, "database-id", "", "Database identifier")
	fs.StringVar(&token, "token", "", "API Token")
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := newClient(token)
	if err != nil {
		return err
	}

	pages, err := client.GetPages(context.Background(), databaseID, nil, nil)
	if err != nil {
		return err
	}
	for _, page := range pages {
		fmt.Printf("ID: %s %+v\n", page.ID, page)
	}

	return nil
}

func getBlocks(args []string) error {
	pageID := ""
	token := ""
	fs := flag.NewFlagSet("get-blocks", flag.ContinueOnError)
	fs.StringVar(&pageID, "page-id", "", "Page identifier")
	fs.StringVar(&token, "token", "", "API Token")
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := newClient(token)
	if err != nil {
		return err
	}

	blocks, err := client.GetBlocks(context.Background(), pageID)
	if err != nil {
		return err
	}
	for _, block := range blocks {
		fmt.Printf("ID: %s %+v\n", block.ID, block)
	}

	return nil
}

func createPage(args []string) error {
	token := ""
	databaseID := ""
	fs := flag.NewFlagSet("create-page", flag.ContinueOnError)
	fs.StringVar(&token, "token", "", "API Token")
	fs.StringVar(&databaseID, "database-id", "", "Parent database identifier")
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := newClient(token)
	if err != nil {
		return err
	}

	db, err := client.GetDatabase(context.Background(), databaseID)
	if err != nil {
		return err
	}

	newPage, err := notion.NewPage(db, "From sample CLI", nil)
	if err != nil {
		return err
	}
	var key string
	for k, v := range db.Properties {
		if v.Type == "rich_text" {
			key = k
			break
		}
	}
	if key != "" {
		newPage.SetProperty(key, &notion.PropertyData{
			Type: "rich_text",
			RichText: []*notion.RichTextObject{
				{
					Type: "text",
					Text: &notion.Text{
						Content: "Test value",
					},
				},
			},
		})
	}
	page, err := client.CreatePage(context.Background(), newPage)
	if err != nil {
		return err
	}
	fmt.Printf("ID: %s\n", page.ID)

	return nil
}

func updateProperties(args []string) error {
	token := ""
	pageID := ""
	fs := flag.NewFlagSet("create-page", flag.ContinueOnError)
	fs.StringVar(&token, "token", "", "API Token")
	fs.StringVar(&pageID, "page-id", "", "Page identifier")
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := newClient(token)
	if err != nil {
		return err
	}

	page, err := client.GetPage(context.Background(), pageID)
	if err != nil {
		return err
	}

	var key string
	for k, v := range page.Properties {
		if v.Type == "rich_text" {
			key = k
			break
		}
	}
	if key == "" {
		return errors.New("text field can not found")
	}

	properties := map[string]*notion.PropertyData{
		key: {
			Type: "rich_text",
			RichText: []*notion.RichTextObject{
				{
					Type: "text",
					Text: &notion.Text{
						Content: "Update property",
					},
				},
			},
		},
	}
	page, err = client.UpdateProperties(context.Background(), pageID, properties)
	if err != nil {
		return err
	}
	fmt.Printf("ID: %s %v\n", page.ID, page.Properties)

	return nil
}

func appendBlocks(args []string) error {
	token := ""
	pageID := ""
	fs := flag.NewFlagSet("append-blocks", flag.ContinueOnError)
	fs.StringVar(&token, "token", "", "API Token")
	fs.StringVar(&pageID, "page-id", "", "Page identifier")
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, err := newClient(token)
	if err != nil {
		return err
	}

	block, err := client.AppendBlock(context.Background(), pageID, []*notion.Block{
		{
			Meta: &notion.Meta{
				Object: "block",
			},
			Type: "paragraph",
			Paragraph: &notion.Paragraph{
				Text: []*notion.RichTextObject{
					{Type: "text", Text: &notion.Text{Content: "Good"}},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf("ID: %s %v\n", block.ID, block)

	return nil
}

func newClient(token string) (*notion.Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return notion.New(tc, notion.BaseURL)
}
