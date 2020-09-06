package main

import (
   "testing"
)

func TestTxReverse(t *testing.T) {
   tests := []struct {
      t Tx
      w Tx
   } {
      {
         Tx { Hash: "abcd1234" },
         Tx { Hash: "3412cdab" },
      },
      {
         Tx { Hash: "008510" },
         Tx { Hash: "108500" },
      },
   }

   for _, tt := range tests {
      ww := tt.t.Reverse()
      e, _ :=  ww.Equals(tt.t)
      if ! e {
         t.Errorf("origin: %v, got: %v, wanted: %v\n", tt.t, ww, tt.w)
      }
   }
}
