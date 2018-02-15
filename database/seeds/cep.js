const neatCsv = require('neat-csv')
const fs = require('fs')
const path = require('path')
const {promisify} = require('util')

const readFile = promisify(fs.readFile)

const tables = [
  {
    name: 'states',
    file: 'states.csv',
    headers: ['id', 'name', 'abbreviation']
  },
  {
    name: 'cities',
    file: 'cities.csv',
    headers: ['id', 'name', 'state_id']
  },
  {
    name: 'postal_codes',
    file: 'ceps.csv',
    headers: ['cep', 'street', 'neighborhood', 'city_id', 'state_id']
  }
]

async function getContent({file, headers}) {
  const content = await readFile(path.resolve(__dirname, file))
  return await neatCsv(content, {
    headers
  })
}

async function seed(knex) {
  for (const table of tables) {
    await knex(table.name).del()
    const content = await getContent(table)
    await knex(table.name).insert(content)
  }
}

module.exports = {
  seed
}
