exports.up = async function(knex) {
  const exists = await knex.schema.hasTable('postal_codes')
  if (exists) {
    return
  }
  await knex.schema.createTable('postal_codes', table => {
    table.increments().primary()
    table.string('cep', 8).notNullable()
    table.unique(['cep'])
    table.string('neighborhood').notNullable()
    table.string('street').notNullable()
    table
      .integer('state_id')
      .unsigned()
      .notNullable()
      .index()
      .references('id')
      .inTable('states')
    table
      .integer('city_id')
      .unsigned()
      .notNullable()
      .index()
      .references('id')
      .inTable('cities')
  })
}

exports.down = function(knex) {
  return knex.schema.dropTable('postal_codes')
}
