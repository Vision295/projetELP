module Main exposing (main)

import Browser
import Html exposing (Html, div, input, text)
import Html.Attributes exposing (placeholder, value)
import Html.Events exposing (onInput)
import TcTurtleParser exposing (Program, read)
import TcTurtleDrawing exposing (display)
import Svg exposing (Svg)


-- MAIN

main : Program () Model Msg
main =
    Browser.sandbox
        { init = init
        , update = update
        , view = view
        }


-- MODEL

type alias Model =
    { input : String
    , program : Result String Program
    }


init : Model
init =
    { input = "[ Forward 100, Right 90, Forward 100 ]"
    , program = read "[ Forward 100, Right 90, Forward 100 ]"
    }


-- UPDATE

type Msg
    = UpdateInput String


update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateInput newInput ->
            { model
                | input = newInput
                , program = read newInput
            }


-- VIEW

view : Model -> Html Msg
view model =
    div []
        [ input
            [ placeholder "Enter TcTurtle program"
            , value model.input
            , onInput UpdateInput
            ]
            []
        , case model.program of
            Ok program ->
                display program

            Err error ->
                Html.pre [] [ text error ]
        ]