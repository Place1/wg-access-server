import formatDistance from 'date-fns/formatDistance';
import timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import { toDate } from './Api';
import { fromResource, lazyObservable } from 'mobx-utils';
import { toast } from './components/Toast';

export function sleep(seconds: number) {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve();
    }, seconds * 1000);
  });
}

export function lastSeen(timestamp: timestamp_pb.Timestamp.AsObject | undefined): string {
  if (timestamp === undefined) {
    return 'Never';
  }
  return formatDistance(toDate(timestamp), new Date(), {
    addSuffix: true,
  });
}

export function lazy<T>(cb: () => Promise<T>) {
  const resource = lazyObservable<T>(async (sink) => {
    sink(await cb());
  });

  return {
    get current() {
      return resource.current();
    },
    refresh: async () => {
      resource.refresh();
    },
  };
}

export function autorefresh<T>(seconds: number, cb: () => Promise<T>) {
  let running = false;
  let sink: ((next: T) => void) | undefined;

  const resource = fromResource<T>(
    async (s) => {
      sink = s;
      running = true;
      while (running) {
        sink(await cb());
        await sleep(seconds);
      }
    },
    () => {
      running = false;
    },
  );

  return {
    get current() {
      return resource.current();
    },
    refresh: async () => {
      if (sink) {
        sink(await cb());
      }
    },
    dispose: () => {
      resource.dispose();
    },
  };
}

export function setClipboard(text: string) {
  const textarea = document.createElement('textarea');
  textarea.value = text;
  document.body.appendChild(textarea);
  textarea.select();
  document.execCommand('copy');
  document.body.removeChild(textarea);
  toast({
    intent: 'success',
    text: 'Added to clipboard',
  });
}

export interface DownloadOpts {
  filename: string,
  content: string,
}

export function download(opts: DownloadOpts) {
  const anchor = document.createElement('a');
  anchor.href = URL.createObjectURL(new File([opts.content], opts.filename));
  anchor.download = opts.filename;
  anchor.style.display = 'none';
  document.body.appendChild(anchor);
  anchor.click();
  document.body.removeChild(anchor);
}
