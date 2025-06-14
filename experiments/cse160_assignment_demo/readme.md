# Running `demo.py`  
A quick-start guide for launching a grading job on the phone-cluster

---

## 1.  Clone the repository
Run these steps after ssh'ing into automaton
```bash
git clone https://github.com/AJ-RR/Junkyard-CSE237D.git
cd Junkyard-CSE237D/experiments/cse160_assignment_demo
```

## 2. Create/Activate a Virtual Environment
```bash
python3 -m venv ~/kube-env
source ~/kube-env/bin/activate
```

## 3. Ensure all dependencies are installed
```bash
pip install --upgrade pip
pip install kubernetes
```

## 4. Run the .py script
```bash
python3 demo.py
```
This should give the following output:

```
Job 'test-cse160-kubernetes' created successfully in namespace 'demo-may-6'.
```

## 5. ssh into the control plane
```bash
ssh user@10.42.0.<#phone>
```

## 6. Watch the job & view logs for the latency information
```bash
kubectl get pods -n demo-may-6 -o wide
```
Then using the pod name which should look like: `test-cse160-kubernetes-<#uniqueID>`
```bash
kubectl logs <pod-name> -n test-assignments-git
```

## 7. Delete job after finished to avoid future errors
```bash
kubectl delete job test-cse160-kubernetes -n demo-may-6
```
