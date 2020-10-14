// 感觉有很多和链表相关的基础操作，增加一个节点，删除一个节点
class Skiplist {
    // 最多16层
    int MAX_LEVEL = 16;
    // 当前SkipList有多少层
    int currLevel;
    Node head;

    public Skiplist() {
        // 默认层数初始化成1
        currLevel = 1;
        // 头节点的值初始化成-1
        head = new Node(-1);
        head.next_point = new Node[MAX_LEVEL];
    }
    
    public boolean search(int target) {
        Node temp = head;
        // 从顶层开始找
        for (int i = MAX_LEVEL -1; i >= 0; i--) {
            while (temp.next_point[i] != null && temp.next_point[i].val <= target) {
                if (temp.next_point[i].val == target) {
                    // 在i层就已经找到了
                    return true;
                } else {
                    // 第i层的下一个节点，右移
                    temp = temp.next_point[i];
                }
            }
        }
        // 在第一层（最底层，数据最全的那一层找），i写死成0
        if (temp.next_point[0] != null && temp.next_point[0].val == target) {
            return true;
        }
        return false;
    }
    
    public void add(int num) {
        // 随机生成层数，如果变大了，就更新
        int level = randomLevel(0.5);
        if (level > currLevel) {
            currLevel = level;
        }
        Node node = new Node(num);
        // 每一层都有next_point，一共level个
        node.next_point = new Node[currLevel];
        Node[] forward = new Node[currLevel];
        
        Arrays.fill(forward, head);
        Node temp = head;
        
        // 找到新插入元素num的前驱节点
        // 从第 level 层开始往第 1 层找
        for (int i = currLevel -1; i >= 0; i--) {
            // next_point[i] 第i层的下一个节点
            while ( temp.next_point[i] != null && temp.next_point[i].val < num) {
                // 一直右移，直到next_point的值比当前值大（模拟右移的过程）
                temp = temp.next_point[i];
            }
            // 更新每一层num的前驱节点
            forward[i] = temp;
        }
        // 更新连接
        for (int i = 0; i < currLevel; i++) {
            // 更新新插入节点的下一个节点，在链表中forward[i] 和 forward[i].next两个节点间插入node节点
            node.next_point[i] = forward[i].next_point[i];
            forward[i].next_point[i] = node;
        }
    }
    
    public boolean erase(int num) {
        // 每一层都有要删除节点的前驱节点
        Node[] forward = new Node[currLevel];
        Node temp = head;
        for (int i = currLevel-1; i >= 0; i--) {
            while (temp.next_point[i] !=null && temp.next_point[i].val < num) {
                temp = temp.next_point[i];
            }
            forward[i] = temp;
        }
        boolean res = false;
        
        // 从最底下一层开始找要删除的节点
        if (temp.next_point[0] != null && temp.next_point[0].val == num) {
            res = true;
            for (int i = currLevel -1; i >= 0; i--) {
                if (forward[i].next_point[i] != null && forward[i].next_point[i].val == num) {
                    forward[i].next_point[i] = forward[i].next_point[i].next_point[i];
                }
            }
        } 
        
        // 删除孤立的层节点
        while (currLevel > 1 && head.next_point[currLevel -1] == null) {
            currLevel -= 1;
        }
        return res;
    }
    
    // 自己也不知道要生成多少层，用随机数吧，小于MAX_LEVEL就行
    public int randomLevel(double p) {
        int level = 1;
        while (Math.random() < p && level < MAX_LEVEL) {
            level++;
        }
        return level;
    }
    
}


class Node {
    int val;
    // 数组下标表示是第几层的
    Node[] next_point;
    public Node(int val){
        this.val = val;
    }
}

/**
 * Your Skiplist object will be instantiated and called as such:
 * Skiplist obj = new Skiplist();
 * boolean param_1 = obj.search(target);
 * obj.add(num);
 * boolean param_3 = obj.erase(num);
 */