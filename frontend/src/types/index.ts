export interface Tag {
    _id: string;
    title: string;
    color: string;
}

export interface Task {
    id: string;
    title: string;
    description: string;
    priority: string;
    completed: boolean;
    tags: string[];
}

export interface Deadline {
    id: string;
    title: string;
    description: string;
    date: string;
    passed: boolean;
    tasks: {
        id?: string;
        title: string;
    }[];
}