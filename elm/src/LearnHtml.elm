module LearnHtml exposing (..)

import Browser
import Html exposing (Html, button, div, input)
import Html.Attributes exposing (placeholder, value)
import Html.Events exposing (onClick, onInput)
import Parser exposing ((|.), (|=), Parser, Trailing(..))
import Svg
import Svg.Attributes exposing (..)

-- MODEL
type alias Model =
    { inputText : String
    , words : List String
    }

init : Model
init =
    { inputText = ""
    , words = []
    }

-- UPDATE
type Msg
    = UpdateText String
    | ParseWords

update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateText newText ->
            { model | inputText = newText }
            
        ParseWords ->
            { model | words = parseText model.inputText }

-- PARSER
parseText : String -> List String
parseText input =
    String.split " " input
        |> List.filter (not << String.isEmpty)

-- VIEW
view : Model -> Html Msg
view model =
      if model.words == ["circle"] then
          div []
              [ Svg.svg 
                  [ width "120"
                  , height "120"
                  , viewBox "0 0 120 120"
                  ] 
                  [ Svg.rect 
                      [ x "10"
                      , y "10"
                      , width "100"
                      , height "100"
                      , rx "15"
                      , ry "15"
                      , fill "blue"
                      ] 
                      []
                  , Svg.circle 
                      [ cx "50"
                      , cy "50"
                      , r "50"
                      , fill "red"
                      ] 
                      []
                  ]
              , div []
                  [ Html.input
                      [ placeholder "Enter some text"
                      , value model.inputText
                      , onInput UpdateText
                      ] 
                      []
                  , button [ onClick ParseWords ] [ Html.text "Send" ]
                  ]
              , div [] 
                  [ Html.text "Words found: "
                  , Html.text (String.fromInt (List.length model.words))
                  ]
              , div [] 
                  [ Html.text "Words: "
                  , Html.text (String.join ", " model.words)
                  ]
              ]
        else
              div[] 
                  [ Html.input
                      [ placeholder "Enter some text"
                      , value model.inputText
                      , onInput UpdateText
                      ] 
                      []
                  , button [ onClick ParseWords ] [ Html.text "Send" ]
                  
              , div [] 
                  [ Html.text "Words found: "
                  , Html.text (String.fromInt (List.length model.words))
                  ]
              , div [] 
                  [ Html.text "Words: "
                  , Html.text (String.join ", " model.words)
                  ]             
                  ]


-- MAIN
main : Program () Model Msg
main =
    Browser.sandbox
        { init = init
        , update = update
        , view = view
        }