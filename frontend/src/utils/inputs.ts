import type { Deadline, DeadlineInput, Task, TaskInput } from "../types";

/** The editable fields of a task, ready to send back to the API. */
export const toTaskInput = (task: Task): TaskInput => ({
    title: task.title,
    description: task.description,
    priority: task.priority,
    completed: task.completed,
    deadline_ids: task.deadline_ids,
    tag_ids: task.tag_ids,
});

/** The editable fields of a deadline, ready to send back to the API. */
export const toDeadlineInput = (deadline: Deadline): DeadlineInput => ({
    title: deadline.title,
    description: deadline.description,
    date: deadline.date,
    tag_ids: deadline.tag_ids,
});

/** Whether two ID lists contain the same IDs, ignoring order. */
export const sameIds = (a: string[], b: string[]): boolean =>
    a.length === b.length && a.every(id => b.includes(id));
