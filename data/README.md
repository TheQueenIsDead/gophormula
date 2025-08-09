# Data

Historic data is available for free from [livetiming.formula1.com/static](https://livetiming.formula1.com/static)

## Testing

The tests within this project are based on historic data from the 2021 Emilia Romagna Grand Prix Race session.

```shell
# Live Timing
STATIC="https://livetiming.formula1.com/static/"
RACE="2021/2021-04-18_Emilia_Romagna_Grand_Prix/2021-04-18_Race/"
CAR="CarData.z.jsonStream"

curl -o "./${RACE}${CAR}" "${STATIC}${RACE}${CAR}"
```