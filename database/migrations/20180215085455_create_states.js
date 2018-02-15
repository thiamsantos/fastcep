exports.up = async function(knex) {
  const exists = await knex.schema.hasTable('states')
  if (exists) {
    return
  }
  await knex.schema.createTable('states', table => {
    table.increments().primary()
    table.string('name').notNullable()
    table.string('abbreviation', 2).notNullable()
  })
}

exports.down = function(knex) {
  return knex.schema.dropTable('states')
}
