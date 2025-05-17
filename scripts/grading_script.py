import subprocess
import re

NAMESPACE = "test-assignments-git"
POD_PREFIX = "grader-job"

def get_grader_pods():
    result = subprocess.run(
        ["kubectl", "get", "pods", "-n", NAMESPACE, "-o", "name"],
        capture_output=True, text=True
    )
    return [line.strip().split('/')[-1] for line in result.stdout.splitlines() if POD_PREFIX in line]

def extract_time(log_text):
    match = re.search(r'real\s+(\d+)m([\d.]+)s', log_text)
    if match:
        minutes = int(match.group(1))
        seconds = float(match.group(2))
        return minutes * 60 + seconds
    return 0

def main():
    total_time = 0
    pod_times = []

    for pod in get_grader_pods():
        result = subprocess.run(
            ["kubectl", "logs", pod, "-n", NAMESPACE],
            capture_output=True, text=True
        )
        time_sec = extract_time(result.stdout)
        total_time += time_sec
        pod_times.append((pod, time_sec))

    for pod, seconds in pod_times:
        print(f"{pod}: {seconds:.2f} sec")

    print(f"\nTotal grading time: {total_time:.2f} sec ({total_time/60:.2f} min)")

if __name__ == "__main__":
    main()
