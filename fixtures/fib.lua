function main(n)
    if n < 2 then return 1 end
    return main(n - 2) + main(n - 1)
end