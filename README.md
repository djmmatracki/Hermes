
# Projekt Hermes

## Wprowadzenie

Hermes jest projektem stworzonym na potrzeby przedmiotu `Badania Operacyjne 2`. Celem projektu bedzie stworzenie API, ktore na podstawie danych podanych przez uzytkownika w formacie <a href="https://www.openstreetmap.org/">Open Street Map</a> bedzie w stanie pomoc firmie transportowej w rozplanowaniu swojej floty w taki sposob, aby transporty byly wykonywane w sposob optymalny. Aplikacja bedzie napisane w jezyku Go Lang z tego wzgledu kod bedzie zoptymalizowany, aby przetwarzac duze ilosci danych w szybki sposob.

## Opis problemu

Zaluzmy ze jestesmy firma przewozowa, ktora posiada swoja flote pojazdow rozlozona w roznych miejscach w kraju. Dostajemy wiele zlecen przewozu towarow, ktore maja wiele miejsc zaladunku oraz wiele miejsc rozladunku.

## Co podlega optymalizacji

Optymalizacji w problemie podlega dlugosc trasy przebytej przez nasza flote oraz zysk, ktory otrzymamy za dane zlecenie. Aplikacja bedzie w stanie obliczac, ktore zlecenie bedzie sie bardziej nam "oplacalo" na podstawie kilometrow przebytych przez flote.

## Istotne uwarunkowania

Zaladunki i rozladunki musza byc zrealizowane w odpowiednich oknach godzinowych. Chcemy aby nasza flota zluzyla jak najmniej paliwa oraz aby ciezar ladunku na ciezarowce nie przekraczal norm. Trzeba pamietac ze ciezarowki nie wiada w kazde miejsce, moga sie pojawic mosty lub inne blokady. Kolejnym uwarunkowaniem jest towar jaki ciezarowki beda przewozic. Moze sie zdarzyc ze nasza flota bedzie miala za zadanie przewisc owoce, ktore potrzebuja odpowiedniej temperatury w pojezdzie. Nie kazda nasza ciezarowka bedzie miala opcje regulowania temperatury.

## Jakie informacje sa potrzebne

Do realizacji aplikacji beda nam potrzebne:

- Dane z open street map
- Informacje o naszej flocie (koordynaty ciezarowek, rodzaje pojazdow)
- Dane o zleceniu (koordynaty zaladunkow oraz rozladunkow, rodzaj przewozonego towaru, cena za zlecenie)

# Uzycie

Stworz plik `.env` przyklad znajduje sie w pliku `env.example`. Wypelnij go danymi zwizanymi z baza danych. Aby zaladowac baze danych uzyj komendy `go run scripts/insertion.go grater-london-latest.osm.pbf`. Aby uruchomic serwer uzyj komendy `go run cmd/main.go`. Domyslnie serwer zostanie uruchomiony na `localhost:8000`. Przetestuj dzialanie serwer'a wpisujac dany adres w przegladarke.

# Endpoint's

- GET /truck - zwraca wszystkie ciezarowki we flocie
- POST /truck - dodaj ciezarowke do floty


Members:

- Patryk Lyczko
- Rafal Maciasz
- Dominik Matracki


- Oplata za autostrade
- Do ktorej musi zdarzyc na zaladunek
- Blokujemy drogi