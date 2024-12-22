package db

// 1. Fetch all events for a group:

// SELECT *
// FROM events
// WHERE group_id = 123;

// 2. Fetch all events created by a specific user:

// SELECT *
// FROM events
// WHERE created_by = 456;

// 3. Fetch all events for a group along with creator details:

// SELECT e.*, u.username
// FROM events e
// JOIN users u ON e.created_by = u.id
// WHERE e.group_id = 123;

// 4. Fetch all groups a user has created events for:

// SELECT DISTINCT g.*
// FROM groups g
// JOIN events e ON g.group_id = e.group_id
// WHERE e.created_by = 456;
