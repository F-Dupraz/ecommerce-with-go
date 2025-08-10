package model

import (
  "string"
)

type Order struct {
  ID string
  UserID string
  ProductID string
  ProductQty uint
}

