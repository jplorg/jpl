# Snippets for performance testing

```jpl
# 📍 Iter / assignment performance tests

#/*
[][99999] = null

# | [.[] | {a:1,b:2}]
# | .[] = {a:1,b:2}
# | a = {a:1,b:2} | [.[] | a]
# | a = {a:1,b:2} | .[] = a

| length()#, .
#*/


# 📍 Nested function call performance tests:

/*
a = [range(0, 400)]
# a = ([0,1,2,3,4,5,6,7,8,9] | [.[],10+.[],20+.[],30+.[],40+.[],50+.[], 60+.[],70+.[],80+.[],90+.[]] | [.[],100+.[],200+.[],300+.[],400+.[],500+.[],600+.[],700+.[],800+.[],900+.[]])

# | [][a[]] = 1
# | a | reduce(func(s): s+., 0)
*/
```
