package main

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Product struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	CategoryID int64  `json:"category_id"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

var products = []Product{
	{ID: 1, Name: "Baju tidur", CategoryID: 1},
	{ID: 2, Name: "Baju renang", CategoryID: 1},
	{ID: 3, Name: "Kursi kaku 8", CategoryID: 2},
	{ID: 4, Name: "Lampu hias", CategoryID: 2},
	{ID: 5, Name: "Meja 360", CategoryID: 2},
	{ID: 6, Name: "Lampu otomatis", CategoryID: 3},
	{ID: 7, Name: "Panel surya", CategoryID: 3},
	{ID: 8, Name: "Palu medium", CategoryID: 4},
	{ID: 9, Name: "Gergaji 2 sisi", CategoryID: 4},
	{ID: 10, Name: "Gerinda ringan", CategoryID: 4},
}

var categories = []Category{
	{ID: 1, Name: "Pakaian"},
	{ID: 2, Name: "Pelengkapan rumah"},
	{ID: 3, Name: "Elektronik"},
	{ID: 4, Name: "Perkakas"},
}

var productType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Product",
		Description: "Represent product",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"category_id": &graphql.Field{
				Type: graphql.Int,
			},
			"category": &graphql.Field{
				Type: categoryType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					srcProduct, ok := p.Source.(Product)
					if ok {
						categoryID := srcProduct.CategoryID
						for _, product := range products {
							if product.CategoryID == categoryID {
								return product, nil
							}
						}
					}
					return nil, nil
				},
			},
		},
	},
)

var categoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Category",
		Description: "Represent category",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Query",
		Description: "Root query",
		Fields: graphql.Fields{
			"product": &graphql.Field{
				Type:        productType,
				Description: "product detail",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						for _, product := range products {
							if int(product.ID) == id {
								return product, nil
							}
						}
					}
					return nil, nil
				},
			},
			"products": &graphql.Field{
				Type:        graphql.NewList(productType),
				Description: "product list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return products, nil
				},
			},
			"category": &graphql.Field{
				Type:        categoryType,
				Description: "category detail",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						for _, category := range categories {
							if int(category.ID) == id {
								return category, nil
							}
						}
					}
					return nil, nil
				},
			},
			"categories": &graphql.Field{
				Type:        graphql.NewList(productType),
				Description: "category list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return categories, nil
				},
			},
		},
	})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Mutation",
	Description: "Root mutation",
	Fields: graphql.Fields{

		"insertProduct": &graphql.Field{
			Type:        productType,
			Description: "Insert new product",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"category_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				product := Product{
					ID:         int64(len(products)) + 1,
					Name:       params.Args["name"].(string),
					CategoryID: params.Args["category_id"].(int64),
				}
				products = append(products, product)
				return product, nil
			},
		},

		"updateProduct": &graphql.Field{
			Type:        productType,
			Description: "Update product by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"category_id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				name, nameOk := params.Args["name"].(string)
				categoryID, categoryIDOk := params.Args["category_id"].(int64)
				product := Product{}
				for i, p := range products {
					if int64(id) == p.ID {
						if nameOk {
							products[i].Name = name
						}
						if categoryIDOk {
							products[i].CategoryID = categoryID
						}
						product = products[i]
						break
					}
				}
				return product, nil
			},
		},

		"deleteProduct": &graphql.Field{
			Type:        productType,
			Description: "Delete product by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				product := Product{}
				for i, p := range products {
					if int64(id) == p.ID {
						product = products[i]
						products = append(products[:i], products[i+1:]...)
					}
				}
				return product, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func main() {
	// http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
	// 	result := executeQuery(r.URL.Query().Get("query"), schema)
	// 	json.NewEncoder(w).Encode(result)
	// })

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)

	fmt.Println("Server is running on port 4000")
	http.ListenAndServe(":4000", nil)
}
