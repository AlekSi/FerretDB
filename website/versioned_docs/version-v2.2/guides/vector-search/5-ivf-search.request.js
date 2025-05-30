db.runCommand({
  aggregate: 'books',
  pipeline: [
    {
      $search: {
        cosmosSearch: {
          vector: [
            0.030855651944875717, 0.038530610501766205, 0.000790110498201102, 0.06512122601270676, 0.009281659498810768,
            -0.05678277835249901, 0.029056841507554054, 0.0216375719755888, 0.012258200906217098, 0.055315714329481125,
            -0.009759286418557167, 0.06137007102370262
          ],
          path: 'vector',
          k: 2
        },
        returnStoredSource: true
      }
    },
    {
      $project: {
        title: 1,
        author: {
          $first: '$authors.name'
        },
        summary: 1,
        vector: 1
      }
    }
  ],
  cursor: {}
})
