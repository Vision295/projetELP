module Main exposing (main)

import Browser
import Html exposing (Html, div, input, svg, text, button)
import Html.Events exposing (onClick, onInput)
import Svg exposing (Svg)
import TcTurtleParser exposing (read)
import TcTurtleDraw exposing (display)

-- MODEL

type alias Model =
    { input : String
    , parsedProgram : Result String Program
    }

type alias Program =
    List Instruction

type Instruction
    = Forward Float
    | Left Float
    | Right Float
    | Repeat Int (List Instruction)

init : Model
init =
    { input = ""
    , parsedProgram = Err "No program yet"
    }


-- UPDATE

type Msg
    = UpdateInput String
    | ParseInput

update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateInput newInput ->
            { model | input = newInput }

        ParseInput ->
            let
                parsed =
                    read model.input
            in
            { model | parsedProgram = parsed }


-- VIEW

view : Model -> Html Msg
view model =
    div []
        [ div []
            [ input [ onInput UpdateInput, Html.Attributes.placeholder "Enter TcTurtle program" ] []
            , button [ onClick ParseInput ] [ text "Parse & Draw" ]
            ]
        , div []
            [ case model.parsedProgram of
                Err error ->
                    text ("Error: " ++ error)

                Ok program ->
                    svg [ Svg.Attributes.viewBox "0 0 500 500", Svg.Attributes.width "500", Svg.Attributes.height "500" ]
                        (display program)
            ]
        ]


-- MAIN

main =
    Browser.sandbox { init = init, update = update, view = view }