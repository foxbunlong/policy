package queries

var GetRole = `
	SELECT
		id, 
		subject, 
		tenant,
		policy,
		created, 
		updated 
		FROM 
		roles
		:where 
		LIMIT ?, ?
`

var InsertRole = `
	INSERT INTO roles(
		subject, 
		tenant,
		policy
	)
	VALUES(
		?,
		?,
		?
	)
`

var DeleteRole = `
	DELETE FROM roles WHERE tenant = ? AND subject = ? AND policy = ? 
`

var GetRoleSubjects = `
	SELECT subject FROM roles WHERE tenant = ? AND policy = ?
`
var GetSubjectRoles = `
	SELECT policy FROM roles WHERE tenant = ? AND subject = ?
`

var GetRoles = `
	SELECT policy FROM roles WHERE tenant = ?
`
