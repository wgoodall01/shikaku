# Shikaku

Solves games of [Shikaku](https://www.wikiwand.com/en/Shikaku).

## What it is

Shiaku is a puzzle game, similar to Sudoku. You start with a board which looks like this:

```
-- -- -- -- -- -- -- 08
07 -- -- -- -- -- -- --
-- -- -- -- -- -- 02 03
-- 08 -- -- -- 04 -- --
-- 04 -- -- -- -- 08 --
-- -- 03 03 -- -- -- --
06 -- -- -- 02 04 -- 02
-- -- -- -- -- -- -- --
```

The goal is to enclose each number in a rectangle of that number's area. None should overlap, and all squares should be covered. 

Basically, any numbered square will be enclosed by a rectangle with dimensions of any of that number's factor pairs.

## How the solver works

The solver has 2 main parts.

First, for each given square on the board, it will determine all the possible bounding rectangles which don't overlap with anything else. Then, if any blank square is only covered by one rectangle, or any given square only has one possible solution, mark the squares as final. This is repeated until no more squares are marked as final.

Then, if the board still has unknown squares, it will pick a possible solution, mark it as final, and try to solve that board, starting again from step 1. It will try possible solutions until one of them eventually works. If the board has no unknown squares, it is solved, and the algorithm returns. If, after all the possible solutions are tried, the solution won't work.


