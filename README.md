# Wordle
Just an implementation of Wordle

## Requirements

- Game consists of
    - target word (5 letters)
    - guess number
    - guess history
    - 26 letter list (used, available, in word, in word in order)


## Server

- start game
- load dict
- take turn (try 5 letters, validate)
- update board

_future features_

- Database to store past wordles
- Stats


## Client

- Type letters
- Send input to server

- Display
    - letter used
    - letter in word, but wrong spot
    - letter in world, correct spot
    - handle multiple letters
    - show previous answers
    - annotate answers (yellow, green)
