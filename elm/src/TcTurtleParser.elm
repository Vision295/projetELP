module TcTurtleParser exposing (TurtleProgram, Instruction(..), read)

import Parser exposing 
    ( Parser, token, int, spaces
    , succeed, (|=), (|.)
    , sequence, run, Problem(..), DeadEnd
    , lazy
    )


-- TYPES

type Instruction
    = Forward Int
    | Left Int
    | Right Int
    | Repeat Int (List Instruction)


type alias TurtleProgram =
    List Instruction


-- PARSING

read : String -> Result String TurtleProgram
read input =
    case run programParser input of
        Ok program ->
            Ok program
            
        Err deadEnds ->
            Err (String.join "\n" (List.map deadEndToString deadEnds))


deadEndToString : DeadEnd -> String
deadEndToString deadEnd =
    "Parser error at row " ++ String.fromInt deadEnd.row ++ 
    ", col " ++ String.fromInt deadEnd.col


programParser : Parser TurtleProgram
programParser =
    succeed identity
        |. spaces
        |. token "["
        |. spaces
        |= instructionSequence
        |. token "]"


instructionSequence : Parser (List Instruction)
instructionSequence =
    sequence
        { start = ""
        , separator = ","
        , end = ""
        , spaces = spaces
        , item = instructionParser
        , trailing = Parser.Optional
        }


instructionParser : Parser Instruction
instructionParser =
    succeed identity
        |. spaces
        |= oneOf
            [ forwardParser
            , leftParser
            , rightParser
            , repeatParser
            ]
        |. spaces


forwardParser : Parser Instruction
forwardParser =
    succeed Forward
        |. token "Forward"
        |. spaces
        |= int


leftParser : Parser Instruction
leftParser =
    succeed Left
        |. token "Left"
        |. spaces
        |= int


rightParser : Parser Instruction
rightParser =
    succeed Right
        |. token "Right"
        |. spaces
        |= int


repeatParser : Parser Instruction
repeatParser =
    succeed Repeat
        |. token "Repeat"
        |. spaces
        |= int
        |. spaces
        |= lazy (\_ -> 
            succeed identity
                |. token "["
                |. spaces
                |= instructionSequence
                |. spaces
                |. token "]"
        )


oneOf : List (Parser a) -> Parser a
oneOf parsers =
    Parser.oneOf parsers