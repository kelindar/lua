local json = require("json")

function main()
    return {
        ["empty"] = json.array(),
        ["empty_map"] = json.array({}),
        ["array"] = json.array({1, 2, 3}),
        ["table"] = json.array({["apple"] = 5}),
        ["str"] = json.array("hello"),
        ["int"] = json.array(12),
        ["bool"] = json.array(true),
        ["float"] = json.array(12.34),
        ["empties"] = json.array(json.array())
    }
end