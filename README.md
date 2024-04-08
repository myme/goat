# Goat

![Goat](./goat.webp)

## 🐐 Goat wherever you wanna goat

Do you have an irrational fear of goats?<br/>
Do you simply enjoy being in the presence of goats?

It's not easy to know where to start looking (or avoid going) when it comes to goats.
Recent intelligence, however, seem to indicate that goats do have a tendency to
flock around farms and surrounding farmland.

This little program works by having you search for your current address
(candidates sorted for your convenience according to your current approximate
position, determined by IP). The program then proceeds to find registered lots
around the entered address where the probability of finding goats `P(goat)` is
greater than `0` (basically property registered as "Farms" in the Norwegian
official cadastre).

## ⚠️ Disclaimer

This satirical little project serves as a personal learning ground for the Go
programming language. The use-case is highly contrived and synthetic, the code
most likely non-idiomatic, and nobody should use this as a base for ...
anything, really.

As for goats? I'd say I'm mostly indifferent to goats. Although goats have
potential utility in humor (fainting goats, and blabbering goats) I, for the
most part, just enjoy the occasional slice of Norwegian, brown goats cheese.

## 🔧 Usage

```
Usage of goat:

Explore places where the probability of goats is greater than 0.

goat <address> <...address>
  -max-results int
        Maximum number of search results (default 10)
  -no-auto-select
        Don't accept first search result automatically
```

Note: All command line arguments are combined into a single address query.

## 🎬 Example

![Goat on the CLI](./goat-cli.gif)
