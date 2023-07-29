// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces

declare global {
  namespace App {
    interface User {
      id: number;
      email: string;
      token: string;
    }
    interface Todo {
      id: number;
      title: string;
      content: string;
      created: Date | null;
      updated: Date | null;
      createdBy: number;
      updatedBy: number;
      done: boolean;
    }
    // interface Error {}
    interface Locals {
      user: User | null;
    }
    // interface PageData {}
    // interface Platform {}
  }
}

export { };
