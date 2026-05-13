import { create } from "zustand";

const useRunsStore = create((set) => ({
  lastRunAt: null,
  notifyRun: () => set({ lastRunAt: Date.now() }),
}));

export default useRunsStore;