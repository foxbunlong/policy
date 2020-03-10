package queries

var GetPolicy = `
	SELECT
		uuid,
		subject, 
		tenant,
		resource,  
		action,
		effect, 
		active, 
		expired, 
		created, 
		updated 
		FROM 
		policies
		:where 
		ORDER BY created DESC LIMIT ?, ?
`

var InsertPolicy = `
	INSERT INTO policies(
		uuid,
		subject, 
		tenant, 
		resource,
		action, 
		effect, 
		active, 
		expired
	)
	VALUES(
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	)
`

var UpdatePolicy = `
	UPDATE policies SET
		subject = ?, 
		tenant = ?, 
		resource = ?,
		action = ?, 
		effect = ?, 
		active = ?, 
		expired = ?,
		updated = ?
	WHERE
		uuid = ?
`
var DeletePolicy = `
	DELETE FROM policies WHERE uuid = ? 
`

var GetResources = `
	SELECT DISTINCT(resource) FROM policies
`

var GetPolicySubjects = `
	SELECT DISTINCT(subject) FROM policies
`
