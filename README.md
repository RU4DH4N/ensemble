## Example

```json
{
  "project": "Example Project",
  "description": "This is an example project to demonstrate JSON structure.",
  "parameters": [
    {
      "id": "variable_1",
      "type": "numeric",
      "min_value": 0,
      "max_value": 100, 
      "step": 0.1
    },
    {
      "id": "variable_2",
      "type": "string",
      "structure": "{adjective} {noun}"
    },
    {
        "id": "variable_3",
        "type": "boolean"
    },
    {
        "id": "variable_4",
        "type": "constant",
        "value": 42
    }
  ],
  "dictionary": {
    "noun": [
      "cat",
      "dog",
      "car"
    ],
    "adjective": [
      "big",
      "small"
    ]
  },
  "seed": 1234
}
```