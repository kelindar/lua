function main()
    return {
        ["items"] = toArray({})
    }
end


function toArray(val)
    -- this indicates to Go that the value is an empty array
    if #val == 0 then
        return "[e7d47667-b92a-48b5-977a-b3199ab09ff9]"
    end
    return val
end