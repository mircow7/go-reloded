# go-reloaded

Petit programme qui corrige un texte.

    go run . sample.txt result.txt

Ce qu'il fait :
- (hex) et (bin) : convertit le nombre d'avant en décimal
- (up) (low) (cap) : majuscules / minuscules / première lettre en majuscule
  (avec un nombre, par ex. (up, 2), ça touche les 2 mots d'avant)
- remet la ponctuation au bon endroit
- colle les guillemets ' ' aux mots
- a devient an devant une voyelle ou un h

Vérifier : go test ./...
