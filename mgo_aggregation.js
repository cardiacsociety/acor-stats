// The aggregation query...
db.data.aggregate([
    {
        $match: {
            siteState: "SA",
            procType: "device"
        }
    },
    {
        $group: {
            _id: {
                month: {$month: "$procDate"},
                year: {$year: "$procDate"},
            },
            count: {$sum: 1},
            date: {$first: "$procDate"}
        }
    },
    { $sort: { "date": 1 } }
])

