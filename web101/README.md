# Övning

## Förberedelser

För att genomföra övningen rekommenderas det att du har docker och docker compose installerat. https://docs.docker.com/compose/install/

## Uppgift

Ändra filerna index.html och style.css

Du ska skapa en sida som visar en lista av todo punkter

Du kan skapa ett nytt föremål genom att posta ett form till '/todo' med ett text fält med namn 'item'
Du kan tabort ett föremål genom att posta ett form till '/todo/:id', där id är föremålets id

Du kan få ut datan i index.html genom att skriva

{{ range .Items }}
    {{.Item}} // Ger text innehållet av punkten
    {{.ID}} // Ger punktens id
{{ end }}
