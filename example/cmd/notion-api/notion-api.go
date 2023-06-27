package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"go.f110.dev/notion-api/v3"
)

func main() {
	cmd := &cobra.Command{
		Use: "notion-api",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}
	cmd.PersistentFlags().String("token", "", "API Token")

	for _, v := range []func(*cobra.Command){
		getUserCmd,
		listUsersCmd,
		getDatabaseCmd,
		updateDatabaseCmd,
		getPageCmd,
		getPagesCmd,
		getPagePropertyCmd,
		getBlockCmd,
		getBlocksCmd,
		updateBlockCmd,
		deleteBlockCmd,
		createPageCmd,
		updatePropertiesCmd,
		appendBlocksCmd,
		searchCmd,
		createDatabaseCmd,
	} {
		v(cmd)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func getUserCmd(parentCmd *cobra.Command) {
	var userID string
	cmd := &cobra.Command{
		Use: "get-user",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
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
		},
	}
	cmd.Flags().StringVar(&userID, "user-id", "", "User identifier")

	parentCmd.AddCommand(cmd)
}

func listUsersCmd(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use: "list-users",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
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
		},
	}

	parentCmd.AddCommand(cmd)
}

func getDatabaseCmd(parentCmd *cobra.Command) {
	var databaseID string
	cmd := &cobra.Command{
		Use: "get-database",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
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
			for _, v := range database.Properties {
				fmt.Printf("%s: %+v\n", v.ID, v)
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&databaseID, "database-id", "", "Database identifier")

	parentCmd.AddCommand(cmd)
}

func updateDatabaseCmd(parentCmd *cobra.Command) {
	var databaseID string
	cmd := &cobra.Command{
		Use: "update-database",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
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
			s := struct{}{}
			newDB := &notion.Database{
				Meta: &notion.Meta{ID: db.ID},
				Properties: map[string]*notion.PropertyMetadata{
					"Foobar": {Name: "Foobar", RichText: &s},
				},
			}
			db, err = client.UpdateDatabase(context.Background(), newDB)
			if err != nil {
				return err
			}
			for _, v := range db.Properties {
				fmt.Printf("%s: %s\n", v.Name, v.Type)
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&databaseID, "database-id", "", "Database identifier")

	parentCmd.AddCommand(cmd)
}

func getPageCmd(parentCmd *cobra.Command) {
	var pageID string
	cmd := &cobra.Command{
		Use: "get-page",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
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
			fmt.Printf("ID: %s %+v\n", pageID, page)
			for name, v := range page.Properties {
				fmt.Printf("%s: %+v\n", name, v)
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&pageID, "page-id", "", "Page identifier")

	parentCmd.AddCommand(cmd)
}

func getPagesCmd(parentCmd *cobra.Command) {
	var databaseID string
	cmd := &cobra.Command{
		Use: "get-pages",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
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
		},
	}
	cmd.Flags().StringVar(&databaseID, "database-id", "", "Database identifier")

	parentCmd.AddCommand(cmd)
}

func getPagePropertyCmd(parentCmd *cobra.Command) {
	var pageID, propertyID string
	cmd := &cobra.Command{
		Use: "get-page-property",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
				return err
			}
			client, err := newClient(token)
			if err != nil {
				return err
			}

			property, err := client.GetPageProperty(context.Background(), pageID, propertyID)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", property)

			return nil
		},
	}
	cmd.Flags().StringVar(&pageID, "page-id", "", "Page identifier")
	cmd.Flags().StringVar(&propertyID, "property-id", "", "Property identifier")

	parentCmd.AddCommand(cmd)
}

func getBlockCmd(parentCmd *cobra.Command) {
	var blockID string
	cmd := &cobra.Command{
		Use: "get-block",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
				return err
			}
			client, err := newClient(token)
			if err != nil {
				return err
			}

			block, err := client.GetBlock(context.Background(), blockID)
			if err != nil {
				return err
			}

			fmt.Printf("ID: %s %+v\n", blockID, block)

			return nil
		},
	}
	cmd.Flags().StringVar(&blockID, "block-id", "", "Block identifier")

	parentCmd.AddCommand(cmd)
}

func getBlocksCmd(parentCmd *cobra.Command) {
	var pageID string
	cmd := &cobra.Command{
		Use: "get-blocks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
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
		},
	}
	cmd.Flags().StringVar(&pageID, "page-id", "", "Page identifier")

	parentCmd.AddCommand(cmd)
}

func updateBlockCmd(parentCmd *cobra.Command) {
	var blockID string
	cmd := &cobra.Command{
		Use: "update-block",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
				return err
			}
			client, err := newClient(token)
			if err != nil {
				return err
			}

			block, err := client.GetBlock(context.Background(), blockID)
			if err != nil {
				return err
			}

			block.Paragraph.Text[0].Text = &notion.Text{Content: "Updated"}

			block, err = client.UpdateBlock(context.Background(), block)
			if err != nil {
				return err
			}
			fmt.Printf("ID: %s %+v\n", block.ID, block)

			return nil
		},
	}
	cmd.Flags().StringVar(&blockID, "block-id", "", "Block iodentifier")

	parentCmd.AddCommand(cmd)
}

func deleteBlockCmd(parentCmd *cobra.Command) {
	var blockID string
	cmd := &cobra.Command{
		Use: "delete-block",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
				return err
			}
			client, err := newClient(token)
			if err != nil {
				return err
			}

			if err := client.DeleteBlock(context.Background(), blockID); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&blockID, "block-id", "", "Block identifier")

	parentCmd.AddCommand(cmd)
}

func createPageCmd(parentCmd *cobra.Command) {
	var databaseID string
	cmd := &cobra.Command{
		Use: "create-page",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
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
		},
	}
	cmd.Flags().StringVar(&databaseID, "database-id", "", "Parent database identifier")

	parentCmd.AddCommand(cmd)
}

func updatePropertiesCmd(parentCmd *cobra.Command) {
	var pageID string
	cmd := &cobra.Command{
		Use: "update-properties",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
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
		},
	}
	cmd.Flags().StringVar(&pageID, "page-id", "", "Page identifier")

	parentCmd.AddCommand(cmd)
}

func appendBlocksCmd(parentCmd *cobra.Command) {
	var pageID string
	cmd := &cobra.Command{
		Use: "append-blocks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
				return err
			}
			client, err := newClient(token)
			if err != nil {
				return err
			}

			blocks, err := client.AppendBlock(context.Background(), pageID, []*notion.Block{
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
			for _, block := range blocks {
				fmt.Printf("ID: %s %+v\n", block.ID, block)
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&pageID, "page-id", "", "Page identifier")

	parentCmd.AddCommand(cmd)
}

func searchCmd(parentCmd *cobra.Command) {
	var query string
	cmd := &cobra.Command{
		Use: "search",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
				return err
			}
			client, err := newClient(token)
			if err != nil {
				return err
			}

			results, err := client.Search(context.Background(), query, nil)
			if err != nil {
				return err
			}
			for _, v := range results {
				switch obj := v.(type) {
				case *notion.Database:
					fmt.Printf("Database ID: %s %+v\n", obj.ID, obj)
				case *notion.Page:
					fmt.Printf("Page ID: %s %+v\n", obj.ID, obj)
				}
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "Search query")

	parentCmd.AddCommand(cmd)
}

func createDatabaseCmd(parentCmd *cobra.Command) {
	var pageID string
	cmd := &cobra.Command{
		Use: "create-database",
		RunE: func(cmd *cobra.Command, _ []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
				return err
			}
			client, err := newClient(token)
			if err != nil {
				return err
			}

			parent, err := client.GetPage(context.Background(), pageID)
			if err != nil {
				return err
			}

			newDatabase := notion.NewDatabase(parent, "Create database")
			newDatabase.SetProperty("Name", &notion.PropertyMetadata{Title: &notion.RichTextObject{}})
			newDatabase.SetProperty("Test1", &notion.PropertyMetadata{RichText: &struct{}{}})
			database, err := client.CreateDatabase(context.Background(), newDatabase)
			if err != nil {
				return err
			}
			fmt.Printf("ID: %s\n", database.ID)

			return nil
		},
	}
	cmd.Flags().StringVar(&pageID, "page-id", "", "Parent page identifier")

	parentCmd.AddCommand(cmd)
}

func newClient(token string) (*notion.Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return notion.New(tc, notion.BaseURL)
}
