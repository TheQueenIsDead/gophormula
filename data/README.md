# Data

Historic data is available for free from [livetiming.formula1.com/static](https://livetiming.formula1.com/static)

## Testing

The tests within this project are based on historic data from the 2021 Emilia Romagna Grand Prix Race session.

While the files are not included in Git, they can be retrieved easily via use of the `historic` program. 

It expects a URI pointing to a folder for a historic session, one that has an `Index.json`, as the first argument.

The second argument is a path to a directory to place the files. The files will be stored in this directory under 
a path that reflects the URI of the main race. 

```shell
go run ./cmd/historic \
  "https://livetiming.formula1.com/static/2021/2021-04-18_Emilia_Romagna_Grand_Prix/2021-04-18_Race/" \
  $PWD/data
```
