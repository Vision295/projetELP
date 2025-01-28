module LearnHtml exposing (main)

import Browser
import Html exposing (Html, div, input, text)
import Html.Attributes exposing (placeholder, value)
import Html.Events exposing (onInput)
import Svg exposing (..)
import Svg.Attributes exposing (..)

-- MODEL
type alias Model =
    { inputText : String
    }

init : Model
init =
    { inputText = ""
    }

-- UPDATE
type Msg
    = UpdateText String

update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateText newText ->
            { model | inputText = newText }

-- VIEW
view : Model -> Html Msg
view model =
    div []
        [ svg 
            [ width "120"
            , height "120"
            , viewBox "0 0 120 120"
            ] 
            [ rect 
                [ x "10"
                , y "10"
                , width "100"
                , height "100"
                , rx "15"
                , ry "15"
                , fill "blue"
                ] 
                []
            , circle 
                [ cx "50"
                , cy "50"
                , r "50"
                , fill "red"
                ] 
                []
            ]
        , Html.input
            [ placeholder "Enter some text"
            , value model.inputText
            , onInput UpdateText
            ] 
            []
        , Html.div [] [ Html.text model.inputText ]
        ]

-- MAIN
main : Program () Model Msg
main =
    Browser.sandbox
        { init = init
        , update = update
        , view = view
        }
