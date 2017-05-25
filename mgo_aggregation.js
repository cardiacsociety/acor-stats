db.Data.aggregate([
    {
        $match: {
            siteState: "NSW",
            procType: "device"
        }
    },
    {
        $group: {
            _id: {
                month: {$month: "$procDate"},
                year: {$year: "$procDate"},
            },
            count: {$sum: 1}
        }
    }
])

