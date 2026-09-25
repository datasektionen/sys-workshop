# Exercise

## Preperation

To perform the exercise it's recommended to use docker compose. https://docs.docker.com/compose/install/

## Assignment

Change the files `index.html` and `style.css`

You're supposed to create a page that shows a list of todo items.

You can create a new item by submiting a form to `/todo` using the post method.
The form needs to contain an input with the name `item`.

You can remove an item by submiting a form to `/todo/{id}`, where `{id}` is the id of the item.

You can get the list of items in `index.html` by using the following syntax (go html/template)

{{ range .Items }}
    {{.Item}} 
    {{.ID}}
{{ end }}

Where `.Item` is the thing entered and `.ID` is a unique id for each item

## Development

Run with automatic rebuild on change

`docker compose up --watch`
