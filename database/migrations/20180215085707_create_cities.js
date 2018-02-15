exports.up = async function(knex) {
  const exists = await knex.schema.hasTable('cities')
  if (exists) {
    return
  }
  await knex.schema.createTable('cities', table => {
    table.increments().primary()
    table
      .integer('state_id')
      .unsigned()
      .notNullable()
      .index()
      .references('id')
      .inTable('states')
    table.string('name').notNullable()
  })
}

exports.down = function(knex) {
  return knex.schema.dropTable('cities')
}
