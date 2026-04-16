package constants

// SQL: pending task operations.
const (
	SQLInsertPendingTask = `INSERT INTO PendingTask
		(TaskTypeId, TargetPath, WorkingDirectory, SourceCommand, CommandArgs)
		VALUES (?, ?, ?, ?, ?)`

	SQLSelectAllPendingTasks = `SELECT p.Id, p.TaskTypeId, t.Name, p.TargetPath,
		p.WorkingDirectory, p.SourceCommand, p.CommandArgs,
		p.FailureReason, p.CreatedAt, p.UpdatedAt
		FROM PendingTask p JOIN TaskType t ON p.TaskTypeId = t.Id
		ORDER BY p.Id`

	SQLSelectPendingTaskByID = `SELECT p.Id, p.TaskTypeId, t.Name, p.TargetPath,
		p.WorkingDirectory, p.SourceCommand, p.CommandArgs,
		p.FailureReason, p.CreatedAt, p.UpdatedAt
		FROM PendingTask p JOIN TaskType t ON p.TaskTypeId = t.Id
		WHERE p.Id = ?`

	SQLSelectPendingTaskByTypePath = `SELECT p.Id FROM PendingTask p
		WHERE p.TaskTypeId = ? AND p.TargetPath = ?`

	SQLSelectPendingTaskByTypePathCmd = `SELECT p.Id FROM PendingTask p
		WHERE p.TaskTypeId = ? AND p.TargetPath = ? AND p.CommandArgs = ?`

	SQLUpdatePendingTaskFailure = `UPDATE PendingTask
		SET FailureReason = ?, UpdatedAt = CURRENT_TIMESTAMP
		WHERE Id = ?`

	SQLDeletePendingTask = `DELETE FROM PendingTask WHERE Id = ?`
)

// SQL: completed task operations.
const (
	SQLInsertCompletedTask = `INSERT INTO CompletedTask
		(OriginalTaskId, TaskTypeId, TargetPath, WorkingDirectory,
		 SourceCommand, CommandArgs, CreatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	SQLSelectAllCompletedTasks = `SELECT c.Id, c.OriginalTaskId, c.TaskTypeId, t.Name,
		c.TargetPath, c.WorkingDirectory, c.SourceCommand, c.CommandArgs,
		c.CompletedAt, c.CreatedAt
		FROM CompletedTask c JOIN TaskType t ON c.TaskTypeId = t.Id
		ORDER BY c.CompletedAt DESC`
)

// SQL: task type lookup.
const SQLSelectTaskTypeByName = `SELECT Id FROM TaskType WHERE Name = ?`
