module Main exposing (main)

import Browser
import Html exposing (Html, div, h1, text)

-- The model (simple unit type in this case)
type alias Model = ()

-- Initialize the model
init : Model
init = ()

-- The view of the page web
view : Model -> Html msg
view model =
    div []
        [ h1 [] [ text "Bienvenue dans Elm!" ]
        , div [] [ text "Ceci est une page web basique." ]
        ]

-- The update function (does nothing for now)
update : msg -> Model -> Model
update msg model = model

-- Configure the program Elm
main =
    Browser.sandbox { init = init, update = update, view = view }