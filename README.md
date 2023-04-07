# Mini tuto API

Pour l'instant, une seule URL est fonctionnelle :
http://localhost:port/shortest-path?depart=nomville&arrivee&nomville

Ici le port est définie sur :8080 (dans le fichier main.go situé à la racine du projet)

Lorsque l'on envoie une requête à cette URL on reçoit une réponse de la forme suivante :
```json
{
    "distance": 186.8000000000001,
    "villeDepart": "Valence",
    "villeArrivee": "Montpellier",
    "pointsReversed": true,
    "points": [],
    "errCode": 0,
    "errMsg": ""
}
```
Où la distance correspond au plus court chemin entre les deux villes. Et où chaque point est un objet de la forme suivante :
```json
{
    "lat": 43.302, // pour la lattitude
    "lon": 4.32 // pour la longitude
}
```
## Code d'erreur
Pour le moment les codes sont les suivants :
* 0 -> pas de problème
* 1 -> pas de chemin trouvé car les villes n'existent pas
* 2 -> erreur dans l'algorithme (impossible inshallah)
* 3 -> le serveur n'est pas encore prêt à recevoir des requêtes