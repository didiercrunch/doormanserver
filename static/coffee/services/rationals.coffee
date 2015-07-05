define ["app", "angularAMD", "lodash"], (app, angularAMD, _) ->

    app.service "rationals", () ->
        service =
            getNominator: (rat) ->
                if _.contains(rat, "/")
                    return Number(rat.substr(0, rat.indexOf("/")))
                Number(rat)

            getDenominator: (rat) ->
                if _.contains(rat, "/")
                    return Number(rat.substr(rat.indexOf("/") + 1))
                1

            getNomDenom: (rat) ->
                return [service.getNominator(rat), service.getDenominator(rat)]

            gcd: (a, b) ->
                if b == 0
                    return a
                return service.gcd(b, a % b)

            lcm: (a, b) ->
                return a * b / service.gcd(a, b)

            changeDenominator: (rat, newDenom) ->
                [nom, denom] = service.getNomDenom(rat)
                return service.newRat( nom *  newDenom / denom,  newDenom)

            add: (rat1, rat2) ->
                [nom, denom] = service.getNomDenom(rat2)
                [nNom, nDenom] = service.getNomDenom(rat1)
                gcd = service.gcd(nDenom, denom)

                rNom = nom*nDenom + nNom * denom
                rDenom = denom * nDenom
                return service.newRat(rNom/gcd , rDenom/gcd)

            minus: (rat1, rat2) ->
                return service.add(rat1, service.opposite(rat2))

            reduce: (rat) ->
                [nom, denom] = service.getNomDenom(rat)
                gcd = service.gcd(nom, denom)
                return service.newRat(nom/gcd, denom/gcd)

            opposite: (rat) ->
                [nom, denom] = service.getNomDenom(rat)
                return service.newRat(-nom, denom)


            sum: (rats...) ->
                ret = "0/1"
                for rat in rats
                    ret = service.add(ret, rat)
                return service.reduce(ret)

            equal: (rat1, rat2) ->
                [nom1, denom1] = service.getNomDenom(rat1)
                [nom2, denom2] = service.getNomDenom(rat2)
                return nom1 * denom2 == nom2 * denom1

            pp: (x) ->
                [num, denum] = service.getNomDenom(service.reduce(x))
                if denum == 1
                    return num
                return "#{num}/#{denum}"


            newRat: (nom , denom) ->
                if nom * denom >= 0
                    return "#{Math.abs(nom)}/#{Math.abs(denom)}"
                return "#{-Math.abs(nom)}/#{Math.abs(denom)}"

        return service

    app.filter "ratpp", ["rationals", (rationals) ->
        filter = (n) ->
            return rationals.pp(n)
        ]
