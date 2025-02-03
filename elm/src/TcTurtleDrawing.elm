module TcTurtleDrawing exposing (display)

import Svg exposing (Svg, svg, polyline)
import Svg.Attributes exposing (viewBox, points, stroke, strokeWidth, fill)
import TcTurtleParser exposing (TurtleProgram, Instruction(..))


-- TYPES

type alias Point =
    { x : Float
    , y : Float
    }


type alias State =
    { position : Point
    , angle : Float
    , path : List Point
    }


-- DISPLAY

display : TurtleProgram -> Svg msg
display program =
    let
        initialState =
            { position = { x = 250, y = 250 }
            , angle = 0
            , path = [ { x = 250, y = 250 } ]
            }

        finalState =
            interpretProgram program initialState
    in
    svg
        [ viewBox "0 0 500 500"
        ]
        [ polyline
            [ points (pointsToString finalState.path)
            , stroke "black"
            , strokeWidth "2"
            , fill "none"
            ]
            []
        ]


-- INTERPRETATION

interpretProgram : TurtleProgram -> State -> State
interpretProgram instructions state =
    List.foldl interpretInstruction state instructions


interpretInstruction : Instruction -> State -> State
interpretInstruction instruction state =
    case instruction of
        Forward distance ->
            let
                angleInRadians =
                    degrees state.angle
                
                newX =
                    state.position.x + (toFloat distance * cos angleInRadians)
                
                newY =
                    state.position.y + (toFloat distance * sin angleInRadians)
                
                newPosition =
                    { x = newX, y = newY }
            in
            { state
                | position = newPosition
                , path = state.path ++ [ newPosition ]
            }

        Left angle ->
            { state | angle = state.angle - toFloat angle }

        Right angle ->
            { state | angle = state.angle + toFloat angle }

        Repeat count instructions ->
            List.foldl
                (\_ accState -> interpretProgram instructions accState)
                state
                (List.range 1 count)


-- HELPERS

pointsToString : List Point -> String
pointsToString points =
    points
        |> List.map (\p -> String.fromFloat p.x ++ "," ++ String.fromFloat p.y)
        |> String.join " "