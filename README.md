## gitfame

Консольная утилита для подсчёта статистик авторов git репозитория.

```
✗ gitfame --repository=. --extensions='.go,.md' --order-by=lines
Name                   Lines Commits Files
Joe Tsai               12154 92      49
colinnewell            130   1       1
Roger Peppe            59    1       2
A. Ishikawa            36    1       1
Tobias Klauser         33    1       2
178inaba               11    2       4
Kyle Lemons            11    1       1
Dmitri Shuralyov       8     1       2
ferhat elmas           7     1       4
Christian Muehlhaeuser 6     3       4
k.nakada               5     1       3
LMMilewski             5     1       2
Ernest Galbrun         3     1       1
Ross Light             2     1       1
Chris Morrow           1     1       1
Fiisio                 1     1       1
```

### Статистики

* Количество строк
* Количество коммитов
* Количество файлов

Все статистики считаются для состояния репозитория на момент конкретного коммита.

### Флаги

Утилита поддерживает следующий набор флагов:

**--repository** — путь до Git репозитория; по умолчанию текущая директория

**--revision** — указатель на коммит; HEAD по умолчанию

**--order-by** — ключ сортировки результатов; один из `lines` (дефолт), `commits`, `files`.

По умолчанию результаты сортируются по убыванию ключа `(lines, commits, files)`.
При равенстве ключей выше будет автор с лексикографически меньшим именем.
При использовании флага соответствующее поле в ключе перемещается на первое место.

**--use-committer** — булев флаг, заменяющий в расчётах автора (дефолт) на коммиттера

**--format** — формат вывода; один из `tabular` (дефолт), `csv`, `json`, `json-lines`;

`tabular`:
```
Name         Lines Commits Files
Joe Tsai     64    3       2
Ross Light   2     1       1
ferhat elmas 1     1       1
```
Human-readable формат. Для паддинга используется пробел.
см. [text/tabwriter](https://golang.org/pkg/text/tabwriter/).

`csv`:
```
Name,Lines,Commits,Files
Joe Tsai,64,3,2
Ross Light,2,1,1
ferhat elmas,1,1,1
```
[encoding/csv](https://golang.org/pkg/encoding/csv/)

`json`:
```
[{"name":"Joe Tsai","lines":64,"commits":3,"files":2},{"name":"Ross Light","lines":2,"commits":1,"files":1},{"name":"ferhat elmas","lines":1,"commits":1,"files":1}]
```
[encoding/json](https://golang.org/pkg/encoding/json/)

`json-lines`:
```
{"name":"Joe Tsai","lines":64,"commits":3,"files":2}
{"name":"Ross Light","lines":2,"commits":1,"files":1}
{"name":"ferhat elmas","lines":1,"commits":1,"files":1}
```

**--extensions** — список расширений, сужающий список файлов в расчёте; множество ограничений разделяется запятыми, например, `'.go,.md'`

**--languages** — список языков (программирования, разметки и др.), сужающий список файлов в расчёте; множество ограничений разделяется запятыми, например `'go,markdown'`

Принадлежность файла к языку программирования определяется с помощью его расширения.
В [configs/language_extensions.json](configs/language_extensions.json) лежит маппинг.
Неизвестные языки никаких ограничений не накладывают. При их использовании можно написать warning в stderr.

**--exclude** — набор [Glob](https://en.wikipedia.org/wiki/Glob_(programming)) паттернов, исключающих файлы из расчёта, например `'foo/*,bar/*'`

Для работы с Glob'ом в стандартной библиотеке есть [path/filepath](https://golang.org/pkg/path/filepath/).

**--restrict-to** — набор Glob паттернов, исключающий все файлы, не удовлетворяющие ни одному из паттернов набора

### Сборка приложения

Как собрать приложение?
```
(cd cmd && go build .)
```
В `cmd` появится исполняемый файл с именем `gitfame`.

Как собрать приложение и установить его в `GOPATH/bin`?
```
go install ./cmd/...
```

Чтобы вызывать установленный бинарь без указания полного пути, нужно добавить `GOPATH/bin` в `PATH`.
```
export PATH=$GOPATH/bin:$PATH
```

После этого `gitfame` будет доступен всюду.

