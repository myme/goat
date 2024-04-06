# Goat

Have an irrational fear of goats?
Simply greatly enjoy being in their presence?

It's not easy to know where to start looking (or avoid) when it comes to goats.
Recent intelligence, however, seem to indicate that goats do have a tendency to
flock around farms and surrounding farmland.

This little program works by having you search for your current address
(candidates sorted for your convenience according to your current approximate
position, determined by IP). The program then proceeds to find registered lots
around the entered address where the probability of finding goats `P(goat)` is
greater than `0`.

## Usage

```
Usage: goat <address> <...address>
```

Note: All command line arguments are combined into a single address query.

## Example

```
❯ go run . myrvollveien oppegård
1: 0.14: Myrvollveien 2B, 1415 OPPEGÅRD
  59.787143, 10.799357
2: 0.14: Myrvollveien 2A, 1415 OPPEGÅRD
  59.786712, 10.799298
3: 0.14: Myrvollveien 4A, 1415 OPPEGÅRD
  59.786210, 10.798459
4: 0.14: Myrvollveien 4C, 1415 OPPEGÅRD
  59.786213, 10.798796
5: 0.14: Myrvollveien 4B, 1415 OPPEGÅRD
  59.785975, 10.798487
Select an address: 1
Gard Ekornrud
     568.0m pos: 59.792100,10.801690
Gard Sætre
    2804.0m pos: 59.762690,10.811130
Gard Haugbru
    2961.0m pos: 59.760990,10.808770
Bruk Kolbotn trelast
     296.0m pos: 59.789780,10.800110
Bruk Nedre Ekornrud
     798.0m pos: 59.794310,10.798760
Bruk Leirskallen
    2255.0m pos: 59.788550,10.759070
Bruk Rosenholm
    4185.0m pos: 59.824650,10.795440
Bruk Øvre Ekornrud
     562.0m pos: 59.790600,10.806650
Bruk Frydenberg
    1041.0m pos: 59.778070,10.803890
```
