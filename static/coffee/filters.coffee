define ['angularAMD', 'app'], (angular, app) ->
    app.filter 'reverse', ->
        (items) ->
            if !Array.isArray(items)
                return false
            items.slice().reverse()
    return
