# Types

Bach has null, boolean, number, string, array, object, and union types.
Additionally, there is the `Void` type (there are no values of this type) and
the `Any` type (all values are of this type). If you are wondering what the
output type of your program is, you can always compose it with the `type`
function, and you will get a string representation of the output type.

    bach> null type
    "Null"
    bach> false type
    "Bool"
    bach> true type
    "Bool"
    bach> 42 type
    "Num"
    bach> 0.3 type
    "Num"
    bach> "Hello world!" type
    "Str"
    bach> [] type
    "Arr<>"
    bach> {} type
    "Obj<Void>"
    bach> for Any def f Num|Str as 0 ok f type
    "Num|Str"
    bach> for Any def f Any as null ok f type
    "Any"
