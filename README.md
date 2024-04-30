# REST API TASK

## run

make run

## endpoints

- post /record

  {
  "startDate": "2017-01-01",
  "endDate": "2018-01-01",
  "minCount": 0,
  "maxCount": 3000
  }

- post /payload

  {
  "key": "key",
  "value": "value"
  }

- get /payload/?key={key}
